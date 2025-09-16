package gobject

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/go-gst/go-glib/pkg/core/classdata"
)

// #cgo pkg-config: gobject-2.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <glib-object.h>
// extern void _goglib_gobject2_Object_constructed(GObject *);
// extern void _goglib_gobject2_Object_dispose(GObject *);
// extern void _goglib_gobject2_Object_get_property(GObject *, guint, GValue *, GParamSpec *);
// extern void _goglib_gobject2_Object_set_property(GObject *, guint, GValue *, GParamSpec *);
// extern void _goglib_gobject2_Object_finalize(GObject *);
// void _goglib_gobject2_Object_virtual_constructed(void* fnptr, GObject* carg0) {
// 	return ((void (*) (GObject*))(fnptr))(carg0);
// }
// void _goglib_gobject2_Object_virtual_dispose(void* fnptr, GObject* carg0) {
// 	return ((void (*) (GObject*))(fnptr))(carg0);
// }
// void _goglib_gobject2_Object_virtual_get_property(void* fnptr, GObject* carg0, guint id, GValue* value, GParamSpec* pspec) {
// 	return ((void (*) (GObject*, guint, GValue*, GParamSpec*))(fnptr))(carg0, id, value, pspec);
// }
// void _goglib_gobject2_Object_virtual_set_property(void* fnptr, GObject* carg0, guint id, GValue* value, GParamSpec* pspec) {
// 	return ((void (*) (GObject*, guint, GValue*, GParamSpec*))(fnptr))(carg0, id, value, pspec);
// }
// void _goglib_gobject2_Object_virtual_finalize(void* fnptr, GObject* carg0) {
// 	return ((void (*) (GObject*))(fnptr))(carg0);
// }
import "C"

type ObjectClass struct {
	*objectClass
}

// objectClass is the struct that is finalized
type objectClass struct {
	native *C.GObjectClass
}

func UnsafeObjectClassFromGlibBorrow(p unsafe.Pointer) *ObjectClass {
	return &ObjectClass{&objectClass{(*C.GObjectClass)(p)}}
}

func UnsafeObjectClassToGlibNone(o *ObjectClass) unsafe.Pointer {
	return unsafe.Pointer(o.native)
}

// InstallProperties will install the given ParameterSpecs to the object class.
// They will be IDed in the order they are provided. Note that the first
// parameter is ID 1, not 0.
func (o *ObjectClass) InstallProperties(params []*ParamSpec) {
	for idx, prop := range params {
		C.g_object_class_install_property(
			(*C.GObjectClass)(UnsafeObjectClassToGlibNone(o)),
			C.guint(idx+1),
			(*C.GParamSpec)(UnsafeParamSpecToGlibNone(prop)),
		)
	}
}

// UnsafeAddPrivateData registers a private structure of the given size for the class
//
// this should not be called by user code, but only by the generated bindings
func (o *ObjectClass) UnsafeAddPrivateData(size uintptr) {
	// https://docs.gtk.org/gobject/method.TypeClass.add_private.html

	// FIXME: this is deprecated, but the alternative is a macro?
	C.g_type_class_add_private(C.gpointer(o.native), C.gsize(size))
}

type ObjectOverrider[Instance Object] interface {
	// getObjectOverrides retrieves the object overrides from any extending overrider
	getObjectOverrides() ObjectOverrides[Instance]
}

// ObjectOverrides is the struct used to override the default implementation of virtual methods.
// it is generic over the extending instance type.
type ObjectOverrides[Instance Object] struct {
	// A callback function used by the type system to initialize a new instance of a type.
	// This function initializes all instance members and allocates any resources required by it.
	// Initialization of a derived instance involves calling all its parent types instance initializers,
	// so the class member of the instance is altered during its initialization to always point to the class that belongs to the type the current initializer was introduced for.
	//
	// The extended members of instance are guaranteed to have been filled with zeros before this function is called.
	//
	// Note: In GObject terms this is not a virtual function on GObject, but instead a function pointer in the GTypeInfo struct.
	// in go-glib we put it here to be able to supply more type information to the user.
	InstanceInit func(instance Instance)

	// The constructed function is called by g_object_new() as the final step of the object creation process.
	// At the point of the call, all construction properties have been set on the object. The purpose of this
	// call is to allow for object initialisation steps that can only be performed after construction properties
	// have been set.
	//
	// The bindings automatically call the parent class constructed method after this function returns.
	Constructed func(Instance)

	// The dispose function is supposed to drop all references to other objects, but keep the instance otherwise intact,
	// so that client method invocations still work. It may be run multiple times (due to reference loops). Before returning,
	// dispose should chain up to the dispose method of the parent class.
	Dispose func(Instance)

	GetProperty func(instance Instance, id uint, pspec *ParamSpec) any
	SetProperty func(instance Instance, id uint, value any, pspec *ParamSpec)

	// Instance finalization function, should finish the finalization of the instance begun in dispose.
	//
	// This is additionally wrapped by the bindings to clean up the instance data.
	//
	// The bindings automatically call the parent class finalize method after this function returns.
	Finalize func(Instance)
}

func (o ObjectOverrides[Instance]) getObjectOverrides() ObjectOverrides[Instance] {
	return o
}

// UnsafeApplyObjectOverrides applies the overrides to init the gclass by setting the trampoline functions.
// This is used by the bindings internally and only exported for visibility to other bindings code.
func UnsafeApplyObjectOverrides[Instance Object](gclass unsafe.Pointer, overrides ObjectOverrides[Instance]) {
	pclass := (*C.GObjectClass)(gclass)

	if overrides.Constructed != nil {
		pclass.constructed = (*[0]byte)(C._goglib_gobject2_Object_constructed)
		classdata.StoreVirtualMethod(
			unsafe.Pointer(pclass),
			"_goglib_gobject2_Object_constructed",
			func(carg0 *C.GObject) {
				var obj Instance // go GObject subclass

				obj = UnsafeObjectFromGlibBorrow(unsafe.Pointer(carg0)).UnsafeLoadInstanceFromPrivateData().(Instance)

				overrides.Constructed(obj)

				obj.ParentConstructed()
			},
		)
	}

	if overrides.Dispose != nil {
		pclass.dispose = (*[0]byte)(C._goglib_gobject2_Object_dispose)
		classdata.StoreVirtualMethod(
			unsafe.Pointer(pclass),
			"_goglib_gobject2_Object_dispose",
			func(carg0 *C.GObject) {
				var obj Instance // go GObject subclass

				obj = UnsafeObjectFromGlibBorrow(unsafe.Pointer(carg0)).UnsafeLoadInstanceFromPrivateData().(Instance)

				overrides.Dispose(obj)
			},
		)
	}

	if overrides.GetProperty != nil {
		pclass.get_property = (*[0]byte)(C._goglib_gobject2_Object_get_property)
		classdata.StoreVirtualMethod(
			unsafe.Pointer(pclass),
			"_goglib_gobject2_Object_get_property",
			func(carg0 *C.GObject, id C.guint, value *C.GValue, pspec *C.GParamSpec) {
				var obj Instance // go GObject subclass
				var param *ParamSpec

				obj = UnsafeObjectFromGlibBorrow(unsafe.Pointer(carg0)).UnsafeLoadInstanceFromPrivateData().(Instance)
				param = UnsafeParamSpecFromGlibNone(unsafe.Pointer(pspec))

				// adjust id to reverse the adjustment from InstallProperties
				goValue := overrides.GetProperty(obj, uint(id)-1, param)

				v := ValueFromNative(unsafe.Pointer(value))

				// SetGoValue will panic if not compatible, but the panic message is not very helpful,
				// because the parameter and instance are missing there. We still panic, but with a more helpful message.
				if !v.CanHold(goValue) {
					panic(fmt.Sprintf("cannot return value of type %T in GValue of type %s as returned by GetProperty override for %T.%s", goValue, v.Type().Name(), obj, param.Name()))
				}

				v.SetGoValue(goValue)
			},
		)
	}

	if overrides.SetProperty != nil {
		pclass.set_property = (*[0]byte)(C._goglib_gobject2_Object_set_property)
		classdata.StoreVirtualMethod(
			unsafe.Pointer(pclass),
			"_goglib_gobject2_Object_set_property",
			func(carg0 *C.GObject, id C.guint, value *C.GValue, pspec *C.GParamSpec) {
				var obj Instance // go GObject subclass
				var param *ParamSpec

				obj = UnsafeObjectFromGlibBorrow(unsafe.Pointer(carg0)).UnsafeLoadInstanceFromPrivateData().(Instance)
				param = UnsafeParamSpecFromGlibNone(unsafe.Pointer(pspec))

				v := ValueFromNative(unsafe.Pointer(value))

				// adjust id to reverse the adjustment from InstallProperties
				overrides.SetProperty(obj, uint(id)-1, v.GoValue(), param)
			},
		)
	}

	// always set the finalize method, because we must clean up the instance data
	pclass.finalize = (*[0]byte)(C._goglib_gobject2_Object_finalize)
	classdata.StoreVirtualMethod(
		unsafe.Pointer(pclass),
		"_goglib_gobject2_Object_finalize",
		func(carg0 *C.GObject) {
			// object and subclass are the same GObject, but we need to be careful with refs here, so we keep both.
			// we have ref=0 here so not all methods are available

			// don't borrow because we don't need the cast
			obj := wrapObject(unsafe.Pointer(carg0))

			subclass := obj.UnsafeLoadInstanceFromPrivateData().(Instance)

			if overrides.Finalize != nil {
				// call the user's finalize first if set.
				// this allows the user to block the finalization of the instance
				// by blocking this call.
				overrides.Finalize(subclass)
			}

			removeInstanceFromPrivateData(subclass)

			obj.ParentFinalize()

			runtime.KeepAlive(subclass)
			runtime.KeepAlive(obj)
		},
	)
}

func RegisterObjectSubClass[InstanceT Object](
	name string,
	classInit func(class *ObjectClass),
	constructor func() InstanceT,
	overrides ObjectOverrides[InstanceT],
	signals map[string]SignalDefinition,
	interfaceInits ...SubClassInterfaceInit[InstanceT],
) Type {
	return UnsafeRegisterSubClass(
		name,
		classInit,
		constructor,
		overrides,
		signals,
		TypeObject,
		UnsafeObjectClassFromGlibBorrow,
		UnsafeApplyObjectOverrides,
		func(obj *ObjectInstance) Object {
			return obj
		},
		interfaceInits...,
	)
}

// Parent virtual method calls on Object

// ParentConstructed calls the parent's constructed virtual method. This should not be needed in user code,
// as the bindings will call this automatically when creating a new instance of the object.
func (obj *ObjectInstance) ParentConstructed() {
	var carg0 *C.GObject

	parentclass := (*C.GObjectClass)(classdata.PeekParentClass(obj.unsafe()))

	carg0 = (*C.GObject)(obj.unsafe())

	C._goglib_gobject2_Object_virtual_constructed(unsafe.Pointer(parentclass.constructed), carg0)
	runtime.KeepAlive(obj)
}

// TODO: ParentDispose, ParentGetProperty, ParentSetProperty

// ParentFinalize calls the parent's finalize virtual method. This should not be needed in user code,
// as the bindings will call this automatically when finalizing the object.
func (obj *ObjectInstance) ParentFinalize() {
	var carg0 *C.GObject

	parentclass := (*C.GObjectClass)(classdata.PeekParentClass(obj.unsafe()))

	carg0 = (*C.GObject)(obj.unsafe())

	C._goglib_gobject2_Object_virtual_finalize(unsafe.Pointer(parentclass.finalize), carg0)
	runtime.KeepAlive(obj)
}
