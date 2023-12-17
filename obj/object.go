package obj

import "github.com/g3n/engine/math32"

var xAixs = math32.Vector3{X: 1, Y: 0, Z: 0}
var yAixs = math32.Vector3{X: 0, Y: 1, Z: 0}
var zAixs = math32.Vector3{X: 0, Y: 0, Z: 1}

type Object struct {
	Id    uint64
	Pos   math32.Vector3
	Rot   math32.Vector3
	Quat  math32.Quaternion
	Local math32.Matrix4
	World math32.Matrix4

	parent                *Object
	children              []*Object
	localAutoUpdate       bool
	worldAutoUpdate       bool
	matrixWorldNeedUpdate bool
}

func NewObject(id uint64) *Object {
	tmp := &Object{
		localAutoUpdate: true,
		worldAutoUpdate: true,
	}
	return tmp
}

func (o *Object) ApplyMatrix4(mat *math32.Matrix4) {
	if o.localAutoUpdate {
		o.UpdateMatrix()
	}

	o.Local.MultiplyMatrices(mat, &o.Local)
	o.Local.Decompose(&o.Pos, &o.Quat, &math32.Vector3{})
}

func (o *Object) ApplyQuaternion(quat *math32.Quaternion) {
	o.Quat.MultiplyQuaternions(quat, &o.Quat)
}

func (o *Object) SetRotationFromAxisAngle(axis *math32.Vector3, angle float32) {
	o.Quat.SetFromAxisAngle(axis, angle)
}

func (o *Object) SetRotationFromEuler(euler *math32.Vector3) {
	o.Quat.SetFromEuler(euler)
}

func (o *Object) SetRotationFromMatrix(mat *math32.Matrix4) {
	o.Quat.SetFromRotationMatrix(mat)
}

func (o *Object) SetRotationFromQuaternion(q *math32.Quaternion) {
	o.Quat.Copy(q)
}

func (o *Object) RotateOnLocalAxis(axis *math32.Vector3, angle float32) {
	q := &math32.Quaternion{}
	q.SetFromAxisAngle(axis, angle)

	o.Quat.MultiplyQuaternions(q, &o.Quat)
}

func (o *Object) RotateOnWorldAxis(axis *math32.Vector3, angle float32) {
	q := &math32.Quaternion{}
	q.SetFromAxisAngle(axis, angle)

	o.Quat.MultiplyQuaternions(&o.Quat, q)
}

func (o *Object) RotateX(angle float32) {
	o.RotateOnLocalAxis(&xAixs, angle)
}

func (o *Object) RotateY(angle float32) {
	o.RotateOnLocalAxis(&yAixs, angle)
}

func (o *Object) RotateZ(angle float32) {
	o.RotateOnLocalAxis(&zAixs, angle)
}

func (o *Object) TranslateOnAxis(axis *math32.Vector3, distance float32) {
	v := axis.Clone().ApplyQuaternion(&o.Quat)
	o.Pos.Add(v.MultiplyScalar(distance))
}

func (o *Object) TranslateX(distance float32) {
	o.TranslateOnAxis(&xAixs, distance)
}

func (o *Object) TranslateY(distance float32) {
	o.TranslateOnAxis(&yAixs, distance)
}

func (o *Object) TranslateZ(distance float32) {
	o.TranslateOnAxis(&zAixs, distance)
}

func (o *Object) LocalToWorld(v *math32.Vector3) *math32.Vector3 {

	o.UpdateWorldMatrix(true, false)

	_v := v.Clone()
	return _v.ApplyMatrix4(&o.World)
}

func (o *Object) WorldToLocal(v *math32.Vector3) *math32.Vector3 {

	o.UpdateWorldMatrix(true, false)

	w := o.World.Clone()
	w.GetInverse(w)

	_v := v.Clone()
	return _v.ApplyMatrix4(w)
}

func (o *Object) UpdateWorldMatrix(updateParent, updateChildren bool) {

	parent := o.parent
	// updatePatent
	if updateParent && parent != nil && parent.worldAutoUpdate {
		// parent auto update
		parent.UpdateWorldMatrix(true, false)
	}

	if o.localAutoUpdate {
		o.UpdateMatrix()
	}

	if parent == nil {
		o.World.Copy(&o.Local)
	} else {
		o.World.MultiplyMatrices(&parent.World, &o.Local)
	}

	if updateChildren {
		for _, child := range o.children {
			if child.worldAutoUpdate {
				child.UpdateWorldMatrix(false, true)
			}
		}
	}
}

func (o *Object) UpdateMatrixWorld(force bool) {
	parent := o.parent

	if o.localAutoUpdate {
		o.UpdateMatrix()
	}

	if o.matrixWorldNeedUpdate || force {
		if parent == nil {
			o.World.Copy(&o.Local)
		} else {
			o.World.MultiplyMatrices(&parent.World, &o.Local)
		}
		o.matrixWorldNeedUpdate = false
		force = true
	}

	for _, child := range o.children {
		if child.worldAutoUpdate || force {
			child.UpdateMatrixWorld(force)
		}
	}
}

func (o *Object) UpdateMatrix() {
	o.Local.Compose(&o.Pos, &o.Quat, &math32.Vector3{X: 1, Y: 1, Z: 1})
	o.matrixWorldNeedUpdate = true
}

func (o *Object) Add(objs ...*Object) {
	for idx := range objs {
		obj := objs[idx]
		if obj == nil {
			continue
		}

		if obj.parent != nil {
			obj.parent.Remove(obj)
		}

		obj.parent = o
		if o.GetObjectById(obj.Id) == nil {
			o.children = append(o.children, obj)
		}
	}
}

func (o *Object) GetObjectById(id uint64) *Object {
	for idx := range o.children {
		obj := o.children[idx]
		if obj.Id == id {
			return obj
		}
	}
	return nil
}

func (o *Object) Remove(objs ...*Object) {
	for _, obj := range objs {
		if obj == nil {
			continue
		}

		for idx := len(o.children) - 1; idx >= 0; idx-- {
			child := o.children[idx]
			if child == nil {
				o.children = append(o.children[:idx], o.children[idx+1:]...)
				continue
			}

			if child.Id != obj.Id {
				continue
			}
			o.children = append(o.children[:idx], o.children[idx+1:]...)

		}
	}
}

type Scene struct {
	Object
	objs map[uint64]*Object
}

func NewScene() *Scene {
	s := &Scene{
		objs: make(map[uint64]*Object),
	}
	return s
}

func (scene *Scene) NewObject() *Object {
	scene.Id += 1

	obj := NewObject(scene.Id)

	scene.objs[obj.Id] = obj

	return obj
}
