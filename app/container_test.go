/*
 *  Copyright (c) 2022-2024 Mikhail Knyazhev <markus621@yandex.ru>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package app_test

import (
	"context"
	"fmt"
	"testing"

	"go.osspkg.com/goppy/app"
	"go.osspkg.com/goppy/errors"
	"go.osspkg.com/goppy/xc"
	"go.osspkg.com/goppy/xtest"
)

func TestUnit_EmptyDI(t *testing.T) {
	c := app.NewContainer(xc.New())
	xtest.NoError(t, c.Start())
	xtest.NoError(t, c.Stop())
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type SimpleDI1_A struct {
	A string
}

func SimpleDI1_func1(a *SimpleDI1_A) {
	fmt.Println("[func]", a.A)
}

func SimpleDI1_func2(a, b *SimpleDI1_A) {
	fmt.Println("[func]", a.A, b.A)
}

type SimpleDI1_Err struct {
	ER string
}

func (v *SimpleDI1_Err) Error() string {
	return v.ER
}

type SimpleDI1_StructEmpty struct {
}

type SimpleDI1_StructUnsupported struct {
	A float32
}

type SimpleDI1_Struct struct {
	AA *SimpleDI1_A
}

type SimpleDI1_ServiceEmpty struct{}

func (v *SimpleDI1_ServiceEmpty) Up() error   { return nil }
func (v *SimpleDI1_ServiceEmpty) Down() error { return nil }

type SimpleDI1_ServiceEmptyXC struct{}

func (v *SimpleDI1_ServiceEmptyXC) Up(_ xc.Context) error { return nil }
func (v *SimpleDI1_ServiceEmptyXC) Down() error           { return nil }

type SimpleDI1_ServiceEmptyCtx struct{}

func (v *SimpleDI1_ServiceEmptyCtx) Up(_ context.Context) error { return nil }
func (v *SimpleDI1_ServiceEmptyCtx) Down() error                { return nil }

type SimpleDI1_Service struct {
	ErrUp   string
	ErrDown string
}

func (v *SimpleDI1_Service) Up() error {
	if len(v.ErrUp) > 0 {
		return fmt.Errorf(v.ErrUp)
	}
	return nil
}

func (v *SimpleDI1_Service) Down() error {
	if len(v.ErrDown) > 0 {
		return fmt.Errorf(v.ErrDown)
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func TestUnit_SimpleDI1(t *testing.T) {
	c := app.NewContainer(xc.New())
	xtest.NoError(t, c.Register(
		&SimpleDI1_A{A: "field A of struct"},
		func(a *SimpleDI1_A) { fmt.Println(1, a.A) },
		func(a *SimpleDI1_A) { fmt.Println(2, a.A) },
		func() { fmt.Println("empty 1") },
		func() { fmt.Println("empty 2") },
		SimpleDI1_func1,
	))
	xtest.NoError(t, c.Start())
	xtest.Error(t, c.Start())
	xtest.NoError(t, c.Stop())
	xtest.NoError(t, c.Stop())
}

func TestUnit_SimpleDI2(t *testing.T) {
	c := app.NewContainer(xc.New())
	xtest.NoError(t, c.Register(&SimpleDI1_A{A: "field A of struct"}))
	xtest.ErrorContains(t, c.Invoke(func(a *SimpleDI1_A) {
		fmt.Println("Invoke", a.A)
	}), "dependencies are not running yet")
	xtest.NoError(t, c.Start())
	xtest.NoError(t, c.Stop())
}

func TestUnit_SimpleDI3(t *testing.T) {
	c := app.NewContainer(xc.New())
	xtest.NoError(t, c.Start())
	xtest.ErrorContains(t, c.Invoke(func(a *SimpleDI1_A) {
		fmt.Println("Invoke", a.A)
	}), "_test.SimpleDI1_A] not initiated")
	xtest.NoError(t, c.Stop())
}

func TestUnit_DI_Default(t *testing.T) {
	tests := []struct {
		name          string
		register      []interface{}
		invoke        interface{}
		wantErr       bool
		wantErrString string
	}{
		{
			name: "Case1",
			register: []interface{}{
				&SimpleDI1_A{A: "field A of struct"},
				&SimpleDI1_A{A: "field A of struct"},
			},
			wantErr:       true,
			wantErrString: "_test.SimpleDI1_A] already initiated",
		},
		{
			name: "Case2",
			register: []interface{}{
				&SimpleDI1_Err{ER: "111"},
			},
			wantErr:       false,
			wantErrString: "",
		},
		{
			name: "Case3",
			register: []interface{}{
				123,
			},
			wantErr:       true,
			wantErrString: "dependency [int] is not supported",
		},
		{
			name: "Case4",
			register: []interface{}{
				&SimpleDI1_A{A: "field A of struct"},
				SimpleDI1_func2,
			},
			wantErr:       false,
			wantErrString: "",
		},
		{
			name: "Case5",
			register: []interface{}{
				&SimpleDI1_A{A: "field A of struct"},
				SimpleDI1_StructEmpty{},
				SimpleDI1_Struct{},
			},
			wantErr:       false,
			wantErrString: "",
		},
		{
			name: "Case6",
			register: []interface{}{
				func() *SimpleDI1_A { return &SimpleDI1_A{A: "field A of struct 1"} },
				func() *SimpleDI1_A { return &SimpleDI1_A{A: "field A of struct 2"} },
			},
			wantErr:       true,
			wantErrString: "_test.SimpleDI1_A] already initiated",
		},
		{
			name: "Case7",
			register: []interface{}{
				SimpleDI1_func2,
			},
			wantErr:       true,
			wantErrString: "_test.SimpleDI1_A] not initiated",
		},
		{
			name: "Case8",
			register: []interface{}{
				SimpleDI1_Struct{},
			},
			wantErr:       true,
			wantErrString: "_test.SimpleDI1_A] not initiated",
		},
		{
			name:          "Case9",
			register:      []interface{}{},
			invoke:        func(a int) {},
			wantErr:       true,
			wantErrString: "dependency [int] is not supported",
		},
		{
			name:          "Case10",
			register:      []interface{}{},
			invoke:        SimpleDI1_StructUnsupported{},
			wantErr:       true,
			wantErrString: "dependency [float32] is not supported",
		},
		{
			name:          "Case11",
			register:      []interface{}{},
			invoke:        SimpleDI1_Struct{},
			wantErr:       true,
			wantErrString: "_test.SimpleDI1_A] not initiated",
		},
		{
			name:          "Case12",
			register:      []interface{}{},
			invoke:        &SimpleDI1_A{},
			wantErr:       false,
			wantErrString: "",
		},
		{
			name:          "Case13",
			register:      []interface{}{},
			invoke:        func() error { return fmt.Errorf("start fail") },
			wantErr:       true,
			wantErrString: "start fail",
		},
		{
			name: "Case14",
			register: []interface{}{
				&SimpleDI1_Service{},
			},
			wantErr:       false,
			wantErrString: "",
		},
		{
			name: "Case15",
			register: []interface{}{
				&SimpleDI1_Service{ErrUp: "is up err"},
			},
			wantErr:       true,
			wantErrString: "is up err",
		},
		{
			name: "Case16",
			register: []interface{}{
				&SimpleDI1_Service{ErrDown: "is down err"},
			},
			wantErr:       true,
			wantErrString: "is down err",
		},
		{
			name: "Case17",
			register: []interface{}{
				func() *SimpleDI1_Service {
					return &SimpleDI1_Service{ErrUp: "is up err"}
				},
			},
			wantErr:       true,
			wantErrString: "is up err",
		},
		{
			name: "Case18",
			register: []interface{}{
				func() (int, error) {
					return 0, nil
				},
			},
			wantErr:       true,
			wantErrString: "dependency [int] is not supported",
		},
		{
			name: "Case19",
			register: []interface{}{
				func() (error, int) {
					return nil, 0
				},
			},
			wantErr:       true,
			wantErrString: "dependency [int] is not supported",
		},
		{
			name:     "Case20",
			register: []interface{}{},
			invoke: func() (int, error) {
				return 0, nil
			},
			wantErr:       false,
			wantErrString: "",
		},
		{
			name:     "Case21",
			register: []interface{}{},
			invoke: func(_ error) int {
				return 0
			},
			wantErr:       true,
			wantErrString: "dependency [error] is not supported",
		},
		{
			name: "Case22",
			register: []interface{}{
				&SimpleDI1_ServiceEmpty{},
				&SimpleDI1_ServiceEmptyXC{},
				&SimpleDI1_ServiceEmptyCtx{},
				&SimpleDI1_Service{},
			},
			wantErr:       false,
			wantErrString: "",
		},
		{
			name: "Case23",
			register: []interface{}{
				func() *SimpleDI1_Service {
					return &SimpleDI1_Service{}
				},
				func(_ *SimpleDI1_Service) error {
					return nil
				},
			},
			wantErr:       false,
			wantErrString: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := app.NewContainer(xc.New())
			errs := errors.Wrap(
				c.Register(tt.register...),
				c.Start(),
				func() error {
					if tt.invoke != nil {
						return c.Invoke(tt.invoke)
					}
					return nil
				}(),
				c.Stop(),
			)
			if tt.wantErr {
				xtest.ErrorContains(t, errs, tt.wantErrString)
			} else {
				xtest.NoError(t, errs)
			}
		})
	}
}
