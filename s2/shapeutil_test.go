/*
Copyright 2017 Google Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package s2

import (
	"testing"

	"github.com/golang/geo/s1"
)

// This file will contain a number of Shape utility types used in different
// parts of testing.
//
//  - edgeVectorShape: represents an arbitrary collection of edges.
//  TODO(roberts): Add remaining testing types here.

// Shape interface enforcement
var (
	_ Shape = (*edgeVectorShape)(nil)
)

// edgeVectorShape is a Shape representing an arbitrary set of edges. It
// is used for testing, but it can also be useful if you have, say, a
// collection of polylines and don't care about memory efficiency (since
// this type would store most of the vertices twice).
type edgeVectorShape struct {
	edges []Edge
}

// edgeVectorShapeFromPoints returns an edgeVectorShape of length 1 from the given points.
func edgeVectorShapeFromPoints(a, b Point) *edgeVectorShape {
	e := &edgeVectorShape{
		edges: []Edge{
			Edge{a, b},
		},
	}
	return e
}

// Add adds the given edge to the shape.
func (e *edgeVectorShape) Add(a, b Point) {
	e.edges = append(e.edges, Edge{a, b})
}
func (e *edgeVectorShape) NumEdges() int                          { return len(e.edges) }
func (e *edgeVectorShape) Edge(id int) Edge                       { return e.edges[id] }
func (e *edgeVectorShape) HasInterior() bool                      { return false }
func (e *edgeVectorShape) ReferencePoint() ReferencePoint         { return OriginReferencePoint(false) }
func (e *edgeVectorShape) NumChains() int                         { return len(e.edges) }
func (e *edgeVectorShape) Chain(chainID int) Chain                { return Chain{chainID, 1} }
func (e *edgeVectorShape) ChainEdge(chainID, offset int) Edge     { return e.edges[chainID] }
func (e *edgeVectorShape) ChainPosition(edgeID int) ChainPosition { return ChainPosition{edgeID, 0} }
func (e *edgeVectorShape) dimension() dimension                   { return polylineGeometry }

func TestEdgeVectorShapeSingletonConstructor(t *testing.T) {
	a := PointFromCoords(1, 0, 0)
	b := PointFromCoords(0, 1, 0)

	var shape Shape = edgeVectorShapeFromPoints(a, b)
	if shape.NumEdges() != 1 {
		t.Errorf("shape created from one edge should only have one edge, got %v", shape.NumEdges())
	}
	if shape.NumChains() != 1 {
		t.Errorf("should only have one edge got %v", shape.NumChains())
	}
	edge := shape.Edge(0)

	if edge.V0 != a {
		t.Errorf("vertex 0 of the edge should be the same as was used to create it. got %v, want %v", edge.V0, a)
	}
	if edge.V1 != b {
		t.Errorf("vertex 1 of the edge should be the same as was used to create it. got %v, want %v", edge.V1, b)
	}
}

// TODO(roberts): Uncomment once LaxLoop and LaxPolygon are added to here.
/*
func TestShapeutilContainsBruteForceNoInterior(t *testing.T) {
	// Defines a polyline that almost entirely encloses the point 0:0.
	polyline := makeLaxPolyline("0:0, 0:1, 1:-1, -1:-1, -1e9:1")
	if containsBruteForce(polyline, parsePoint("0:0")) {
		t.Errorf("containsBruteForce(%v, %v) = true, want false")
	}
}

func TestShapeutilContainsBruteForceContainsReferencePoint(t *testing.T) {
	// Checks that containsBruteForce agrees with ReferencePoint.
	polygon := makeLaxPolygon("0:0, 0:1, 1:-1, -1:-1, -1e9:1")
	ref, _ := polygon.ReferencePoint()
	if got := containsBruteForce(polygon, ref.point); got != ref.contained {
		t.Errorf("containsBruteForce(%v, %v) = %v, want %v", polygon, ref.Point, got, ref.contained)
	}
}
*/

func TestShapeutilContainsBruteForceConsistentWithLoop(t *testing.T) {
	// Checks that containsBruteForce agrees with Loop Contains.
	loop := RegularLoop(parsePoint("89:-179"), s1.Angle(10)*s1.Degree, 100)
	for i := 0; i < loop.NumVertices(); i++ {
		if got, want := loop.ContainsPoint(loop.Vertex(i)),
			containsBruteForce(loop, loop.Vertex(i)); got != want {
			t.Errorf("loop.ContainsPoint(%v) = %v, containsBruteForce(shape, %v) = %v, should be the same", loop.Vertex(i), got, loop.Vertex(i), want)
		}
	}
}
