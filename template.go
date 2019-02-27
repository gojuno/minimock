package minimock

const (
	// HeaderTemplate is used to generate package clause and go:generate instruction
	HeaderTemplate = `
		package {{$.Package.Name}}

		// DO NOT EDIT!
		// The code below was generated with http://github.com/gojuno/minimock ({{$.Options.HeaderVars.Version}})

		{{if $.Options.HeaderVars.GenerateInstruction}}
		//go:generate minimock -i {{$.SourcePackage.PkgPath}}.{{$.Options.InterfaceName}} -o {{$.Options.OutputFile}}
		{{end}}

		import (
			"sync/atomic"
			"time"

			{{range $import := $.Options.Imports}}
				{{if not (in $import "\"time\"" "\"sync/atomic\"" "minimock \"github.com/gojuno/minimock/pkg\"")}}{{$import}}{{end}}
			{{end}}

			minimock "github.com/gojuno/minimock/pkg"
		)
	`

	// BodyTemplate is used to generate mock body
	BodyTemplate = `
		{{ $mock := (title (printf "%sMock" $.Interface.Name)) }}

		// {{$mock}} implements {{$.Interface.Type}}
		type {{$mock}} struct {
			t minimock.Tester
			{{ range $method := $.Interface.Methods }}
				func{{$method.Name}} func{{ $method.Signature }}
				after{{$method.Name}}Counter uint64
				before{{$method.Name}}Counter uint64
				{{$method.Name}}Mock m{{$mock}}{{$method.Name}}
			{{ end }}
		}

		// New{{$mock}} returns a mock for {{$.Interface.Type}}
		func New{{$mock}}(t minimock.Tester) *{{$mock}} {
			m := &{{$mock}}{t: t}
			if controller, ok := t.(minimock.MockController); ok {
				controller.RegisterMocker(m)
			}
			{{ range $method := $.Interface.Methods }}m.{{$method.Name}}Mock = m{{$mock}}{{$method.Name}}{mock: m}
			{{ end }}
			return m
		}

		{{ range $method := $.Interface.Methods }}
			type m{{$mock}}{{$method.Name}} struct {
				mock              *{{$mock}}
				defaultExpectation   *{{$mock}}{{$method.Name}}Expectation
				expectations []*{{$mock}}{{$method.Name}}Expectation
			}

			// {{$mock}}{{$method.Name}}Expectation specifies expectation struct of the {{$.Interface.Name}}.{{$method.Name}}
			type {{$mock}}{{$method.Name}}Expectation struct {
				mock *{{$mock}}
				{{ if $method.HasParams }}  params *{{$mock}}{{$method.Name}}Params  {{end}}
				{{ if $method.HasResults }} results *{{$mock}}{{$method.Name}}Results {{end}}
				Counter uint64
			}

			{{if $method.HasParams }}
				// {{$mock}}{{$method.Name}}Params contains parameters of the {{$.Interface.Name}}.{{$method.Name}}
				type {{$mock}}{{$method.Name}}Params {{$method.ParamsStruct}}
			{{end}}

			{{if $method.HasResults }}
				// {{$mock}}{{$method.Name}}Results contains results of the {{$.Interface.Name}}.{{$method.Name}}
				type {{$mock}}{{$method.Name}}Results {{$method.ResultsStruct}}
			{{end}}

			// Expect sets up expected params for {{$.Interface.Name}}.{{$method.Name}}
			func (m *m{{$mock}}{{$method.Name}}) Expect({{$method.Params}}) *m{{$mock}}{{$method.Name}} {
				if m.mock.func{{$method.Name}} != nil {
					m.mock.t.Fatalf("{{$mock}}.{{$method.Name}} mock is already set by Set")
				}

				if m.defaultExpectation == nil {
					m.defaultExpectation = &{{$mock}}{{$method.Name}}Expectation{}
				}

				{{if $method.HasParams }}
					m.defaultExpectation.params = &{{$mock}}{{$method.Name}}Params{ {{ $method.ParamsNames }} }
					for _, e := range m.expectations {
						if minimock.Equal(e.params, m.defaultExpectation.params) {
							m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
						}
					}
				{{end}}
				return m
			}

			// Return sets up results that will be returned by {{$.Interface.Name}}.{{$method.Name}}
			func (m *m{{$mock}}{{$method.Name}}) Return({{$method.Results}}) *{{$mock}} {
				if m.mock.func{{$method.Name}} != nil {
					m.mock.t.Fatalf("{{$mock}}.{{$method.Name}} mock is already set by Set")
				}

				if m.defaultExpectation == nil {
					m.defaultExpectation = &{{$mock}}{{$method.Name}}Expectation{mock: m.mock}
				}
				{{if $method.HasResults }} m.defaultExpectation.results = &{{$mock}}{{$method.Name}}Results{ {{ $method.ResultsNames }} } {{end}}
				return m.mock
			}

			//Set uses given function f to mock the {{$.Interface.Name}}.{{$method.Name}} method
			func (m *m{{$mock}}{{$method.Name}}) Set(f func{{$method.Signature}}) *{{$mock}}{
				if m.defaultExpectation != nil {
					m.mock.t.Fatalf("Default expectation is already set for the {{$.Interface.Name}}.{{$method.Name}} method")
				}

				if len(m.expectations) > 0 {
					m.mock.t.Fatalf("Some expectations are already set for the {{$.Interface.Name}}.{{$method.Name}} method")
				}

				m.mock.func{{$method.Name}}= f
				return m.mock
			}

			{{if (and $method.HasParams $method.HasResults)}}
				// When sets expectation for the {{$.Interface.Name}}.{{$method.Name}} which will trigger the result defined by the following
				// Then helper
				func (m *m{{$mock}}{{$method.Name}}) When({{$method.Params}}) *{{$mock}}{{$method.Name}}Expectation {
					if m.mock.func{{$method.Name}} != nil {
						m.mock.t.Fatalf("{{$mock}}.{{$method.Name}} mock is already set by Set")
					}

					expectation := &{{$mock}}{{$method.Name}}Expectation{
						mock: m.mock,
						params: &{{$mock}}{{$method.Name}}Params{ {{ $method.ParamsNames }} },
					}
					m.expectations = append(m.expectations, expectation)
					return expectation
				}

				// Then sets up {{$.Interface.Name}}.{{$method.Name}} return parameters for the expectation previously defined by the When method
				func (e *{{$mock}}{{$method.Name}}Expectation) Then({{$method.Results}}) *{{$mock}} {
					e.results = &{{$mock}}{{$method.Name}}Results{ {{ $method.ResultsNames }} }
					return e.mock
				}
			{{end}}

			// {{$method.Name}} implements {{$.Interface.Type}}
			func (m *{{$mock}}) {{$method.Declaration}} {
				atomic.AddUint64(&m.before{{$method.Name}}Counter, 1)
				defer atomic.AddUint64(&m.after{{$method.Name}}Counter, 1)

				{{if $method.HasParams}}
					for _, e := range m.{{$method.Name}}Mock.expectations {
						if minimock.Equal(*e.params,  {{$mock}}{{$method.Name}}Params{ {{$method.ParamsNames}} }) {
							atomic.AddUint64(&e.Counter, 1)
							{{$method.ReturnStruct "e.results" -}}
						}
					}
				{{end}}

				if m.{{$method.Name}}Mock.defaultExpectation != nil {
					atomic.AddUint64(&m.{{$method.Name}}Mock.defaultExpectation.Counter, 1)
					{{- if $method.HasParams }}
						want:= m.{{$method.Name}}Mock.defaultExpectation.params
						got:= {{$mock}}{{$method.Name}}Params{ {{$method.ParamsNames}} }
						if want != nil && !minimock.Equal(*want, got) {
							m.t.Errorf("{{$mock}}.{{$method.Name}} got unexpected parameters, want: %#v, got: %#v\n", *want, got)
						}
					{{ end }}
					{{if $method.HasResults }}
						results := m.{{$method.Name}}Mock.defaultExpectation.results
						if results == nil {
							m.t.Fatal("No results are set for the {{$mock}}.{{$method.Name}}")
						}
						{{$method.ReturnStruct "(*results)" -}}
					{{else}}
						return
					{{ end }}
				}
				if m.func{{$method.Name}} != nil {
					{{$method.Pass "m.func"}}
				}
				m.t.Fatalf("Unexpected call to {{$mock}}.{{$method.Name}}.{{range $method.Params}} %v{{end}}", {{ $method.ParamsNames }} )
				{{if $method.HasResults}}return{{end}}
			}

			// {{$method.Name}}AfterCounter returns a count of finished {{$mock}}.{{$method.Name}} invocations
			func (m *{{$mock}}) {{$method.Name}}AfterCounter() uint64 {
				return atomic.LoadUint64(&m.after{{$method.Name}}Counter)
			}

			// {{$method.Name}}BeforeCounter returns a count of {{$mock}}.{{$method.Name}} invocations
			func (m *{{$mock}}) {{$method.Name}}BeforeCounter() uint64 {
				return atomic.LoadUint64(&m.before{{$method.Name}}Counter)
			}

			// Minimock{{$method.Name}}Done returns true if the count of the {{$method.Name}} invocations corresponds
			// the number of defined expectations
			func (m *{{$mock}}) Minimock{{$method.Name}}Done() bool {
				for _, e := range m.{{$method.Name}}Mock.expectations {
					if atomic.LoadUint64(&e.Counter) < 1 {
						return false
					}
				}

				// if default expectation was set then invocations count should be greater than zero
				if m.{{$method.Name}}Mock.defaultExpectation != nil && atomic.LoadUint64(&m.after{{$method.Name}}Counter) < 1 {
					return false
				}
				// if func was set then invocations count should be greater than zero
				if m.func{{$method.Name}} != nil && atomic.LoadUint64(&m.after{{$method.Name}}Counter) < 1  {
					return false
				}
				return true
			}

			// Minimock{{$method.Name}}Inspect logs each unmet expectation
			func (m *{{$mock}}) Minimock{{$method.Name}}Inspect() {
				for _, e := range m.{{$method.Name}}Mock.expectations {
					if atomic.LoadUint64(&e.Counter) < 1 {
						{{- if $method.HasParams}}
							m.t.Errorf("Expected call to {{$mock}}.{{$method.Name}} with params: %#v", *e.params)
						{{else}}
							m.t.Error("Expected call to {{$mock}}.{{$method.Name}}")
						{{end -}}
					}
				}

				// if default expectation was set then invocations count should be greater than zero
				if m.{{$method.Name}}Mock.defaultExpectation != nil && atomic.LoadUint64(&m.after{{$method.Name}}Counter) < 1 {
					{{- if $method.HasParams}}
						m.t.Errorf("Expected call to {{$mock}}.{{$method.Name}} with params: %#v", *m.{{$method.Name}}Mock.defaultExpectation.params)
					{{else}}
						m.t.Error("Expected call to {{$mock}}.{{$method.Name}}")
					{{end -}}
				}
				// if func was set then invocations count should be greater than zero
				if m.func{{$method.Name}} != nil && atomic.LoadUint64(&m.after{{$method.Name}}Counter) < 1  {
					m.t.Error("Expected call to {{$mock}}.{{$method.Name}}")
				}
			}
		{{end}}

		// MinimockFinish checks that all mocked methods have been called the expected number of times
		func (m *{{$mock}}) MinimockFinish() {
			if !m.minimockDone() {
				{{- range $method := $.Interface.Methods }}
					m.Minimock{{$method.Name}}Inspect()
				{{ end -}}
				m.t.FailNow()
			}
		}

		// MinimockWait waits for all mocked methods to be called the expected number of times
		func (m *{{$mock}}) MinimockWait(timeout time.Duration) {
			timeoutCh := time.After(timeout)
			for {
				if m.minimockDone() {
					return
				}
				select {
				case <-timeoutCh:
					m.MinimockFinish()
					return
				case <-time.After(10 * time.Millisecond):
				}
			}
		}

		func (m *{{$mock}}) minimockDone() bool {
			done := true
			return done {{ range $method := $.Interface.Methods }}&&
			m.Minimock{{$method.Name}}Done(){{end -}}
		}
	`
)
