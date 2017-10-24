package decoder

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	unstructuredconversion "k8s.io/apimachinery/pkg/conversion/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"

	admissionregistrationv1alpha1 "k8s.io/api/admissionregistration/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	authenticationv1 "k8s.io/api/authentication/v1"
	authenticationv1beta1 "k8s.io/api/authentication/v1beta1"
	authorizationv1 "k8s.io/api/authorization/v1"
	authorizationv1beta1 "k8s.io/api/authorization/v1beta1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	autoscalingv2beta1 "k8s.io/api/autoscaling/v2beta1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	batchv2alpha1 "k8s.io/api/batch/v2alpha1"
	certificatesv1beta1 "k8s.io/api/certificates/v1beta1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	rbacv1alpha1 "k8s.io/api/rbac/v1alpha1"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
	schedulingv1alpha1 "k8s.io/api/scheduling/v1alpha1"
	settingsv1alpha1 "k8s.io/api/settings/v1alpha1"
	storagev1 "k8s.io/api/storage/v1"
	storagev1beta1 "k8s.io/api/storage/v1beta1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	// Scheme knows about all kubernetes types
	Scheme         = runtime.NewScheme()
	Codecs         = serializer.NewCodecFactory(Scheme)
	ParameterCodec = runtime.NewParameterCodec(Scheme)
)

func init() {

	v1.AddToGroupVersion(Scheme, schema.GroupVersion{Version: "v1"})
	AddToScheme(Scheme)
}

func AddToScheme(scheme *runtime.Scheme) {
	admissionregistrationv1alpha1.AddToScheme(scheme)
	appsv1beta1.AddToScheme(scheme)
	appsv1beta2.AddToScheme(scheme)
	appsv1.AddToScheme(scheme)
	authenticationv1.AddToScheme(scheme)
	authenticationv1beta1.AddToScheme(scheme)
	authorizationv1.AddToScheme(scheme)
	authorizationv1beta1.AddToScheme(scheme)
	autoscalingv1.AddToScheme(scheme)
	autoscalingv2beta1.AddToScheme(scheme)
	batchv1.AddToScheme(scheme)
	batchv1beta1.AddToScheme(scheme)
	batchv2alpha1.AddToScheme(scheme)
	certificatesv1beta1.AddToScheme(scheme)
	corev1.AddToScheme(scheme)
	extensionsv1beta1.AddToScheme(scheme)
	networkingv1.AddToScheme(scheme)
	policyv1beta1.AddToScheme(scheme)
	rbacv1.AddToScheme(scheme)
	rbacv1beta1.AddToScheme(scheme)
	rbacv1alpha1.AddToScheme(scheme)
	schedulingv1alpha1.AddToScheme(scheme)
	settingsv1alpha1.AddToScheme(scheme)
	storagev1beta1.AddToScheme(scheme)
	storagev1.AddToScheme(scheme)

}

type streamObject struct {
	obj runtime.Object
	err error
}

func FileToKubeObj(fileName string) ([]runtime.Object, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	d := streamingDecoder(f)
	// can decode a list of upto 100 objects
	out := make(chan *streamObject, 100)
	obj := &unstructured.Unstructured{}
	err = d.Decode(obj)
	switch {
	case err == io.EOF:
	case err != nil:
		return nil, err
	default:
		out <- &streamObject{obj: obj}
	}
	results := flatten(out)
	results = typecast(results, Scheme)
	results = convertinternal(results, Scheme)

	close(out)

	returnObj := []runtime.Object{}
	for x := range results {
		returnObj = append(returnObj, x.obj)
		if x.err != nil {
			return nil, x.err
		}
	}
	return returnObj, nil
}

// decoder can decode streaming json, yaml docs, single json objects, single yaml objects
type decoder interface {
	Decode(into interface{}) error
}

func streamingDecoder(r io.ReadCloser) decoder {
	buffer := bufio.NewReaderSize(r, 1024)
	b, _ := buffer.Peek(1)
	if string(b) == "{" {
		return json.NewDecoder(buffer)
	} else {
		return yaml.NewYAMLToJSONDecoder(buffer)
	}
}

func stream(sources []io.ReadCloser) <-chan *streamObject {
	out := make(chan *streamObject)

	wg := &sync.WaitGroup{}
	for i := range sources {
		wg.Add(1)
		go func(r io.ReadCloser) {
			defer wg.Done()
			defer r.Close()
			d := streamingDecoder(r)
			for {
				obj := &unstructured.Unstructured{}
				err := d.Decode(obj)
				switch {
				case err == io.EOF:
					return
				case err != nil:
					out <- &streamObject{err: err}
				default:
					out <- &streamObject{obj: obj}
				}
			}
		}(sources[i])
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func flatten(in <-chan *streamObject) <-chan *streamObject {
	out := make(chan *streamObject)

	v1List := v1.SchemeGroupVersion.WithKind("List")

	go func() {
		defer close(out)
		for result := range in {
			if result.err != nil {
				out <- result
				continue
			}

			if result.obj.GetObjectKind().GroupVersionKind() != v1List {
				out <- result
				continue
			}

			data, err := json.Marshal(result.obj)
			if err != nil {
				out <- &streamObject{err: err}
				continue
			}

			list := &unstructured.UnstructuredList{}
			if err := list.UnmarshalJSON(data); err != nil {
				out <- &streamObject{err: err}
				continue
			}

			for _, item := range list.Items {
				out <- &streamObject{obj: &item}
			}
		}
	}()
	return out
}

func typecast(in <-chan *streamObject, creator runtime.ObjectCreater) <-chan *streamObject {
	out := make(chan *streamObject)

	go func() {
		defer close(out)
		for result := range in {
			if result.err != nil {
				out <- result
				continue
			}

			typed, err := creator.New(result.obj.GetObjectKind().GroupVersionKind())
			if err != nil {
				out <- &streamObject{err: err}
				continue
			}

			unstructuredObject, ok := result.obj.(*unstructured.Unstructured)
			if !ok {
				out <- &streamObject{err: fmt.Errorf("expected *unstructured.Unstructured, got %T", result.obj)}
			}

			if err := unstructuredconversion.DefaultConverter.FromUnstructured(unstructuredObject.Object, typed); err != nil {
				out <- &streamObject{err: err}
				continue
			}

			out <- &streamObject{obj: typed}
		}
	}()
	return out
}

func convertinternal(in <-chan *streamObject, convertor runtime.ObjectConvertor) <-chan *streamObject {
	out := make(chan *streamObject)

	go func() {
		defer close(out)
		for result := range in {
			if result.err != nil {
				out <- result
				continue
			}

			gv := result.obj.GetObjectKind().GroupVersionKind().GroupVersion()
			if gv.Version == "" || gv.Version == runtime.APIVersionInternal {
				out <- result
				continue
			}

			//gv.Version = runtime.APIVersionInternal
			converted, err := convertor.ConvertToVersion(result.obj, gv)
			if err != nil {
				fmt.Println("xyz")
				out <- &streamObject{err: err}
				continue
			}

			out <- &streamObject{obj: converted}
		}
	}()
	return out
}
