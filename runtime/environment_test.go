package runtime_test

import (
	"fmt"
	"github.com/go-sre/host/runtime"
)

func ExampleDevEnv() {
	fmt.Printf("test: IsDevEnv() -> %v\n", runtime.IsDevEnv())

	runtime.SetRuntimeEnv("dev")
	fmt.Printf("test: IsDevEnv(dev) -> %v\n", runtime.IsDevEnv())

	runtime.SetRuntimeEnv("devrrr")
	fmt.Printf("test: IsDevEnv(devrrr) -> %v\n", runtime.IsDevEnv())

	// Output:
	//test: IsDevEnv() -> true
	//test: IsDevEnv(dev) -> true
	//test: IsDevEnv(devrrr) -> false

}

/*
func ExampleProdEnv() {
	fmt.Println(resource.IsProdEnv())
	os.Setenv(resource.RuntimeEnvKey, "prod")
	fmt.Println(resource.IsProdEnv())
	os.Setenv(resource.RuntimeEnvKey, "production")
	fmt.Println(resource.IsProdEnv())

	// Output:
	// false
	// true
	// false
}

func ExampleReviewEnv() {
	fmt.Println(resource.IsReviewEnv())
	os.Setenv(resource.RuntimeEnvKey, "review")
	fmt.Println(resource.IsReviewEnv())
	os.Setenv(resource.RuntimeEnvKey, "revvrrr")
	fmt.Println(resource.IsReviewEnv())

	// Output:
	// false
	// true
	// false
}

func ExampleStageEnv() {
	fmt.Println(resource.IsStageEnv())
	os.Setenv(resource.RuntimeEnvKey, "stage")
	fmt.Println(resource.IsStageEnv())
	os.Setenv(resource.RuntimeEnvKey, "")
	fmt.Println(resource.IsStageEnv())

	// Output:
	// false
	// true
	// false
}

func ExampleTestEnv() {
	fmt.Println(resource.IsTestEnv())
	os.Setenv(resource.RuntimeEnvKey, "test")
	fmt.Println(resource.IsTestEnv())
	os.Setenv(resource.RuntimeEnvKey, "atvrrr")
	fmt.Println(resource.IsTestEnv())

	// Output:
	// false
	// true
	// false
}


*/

func ExampleExpand() {
	runtime.SetRuntimeEnv("dev")

	t := "file_{env}.txt"
	fmt.Printf("test: EnvExpansion(%v) -> %v\n", t, runtime.EnvExpansion(t))

	runtime.SetRuntimeEnv("prod")
	fmt.Printf("test: EnvExpansion(%v) -> %v\n", t, runtime.EnvExpansion(t))

	t = "file_{env_var_invalid}.txt"
	fmt.Printf("test: EnvExpansion(%v) -> %v\n", t, runtime.EnvExpansion(t))

	// Output:
	//test: EnvExpansion(file_{env}.txt) -> file_dev.txt
	//test: EnvExpansion(file_{env}.txt) -> file_prod.txt
	//test: EnvExpansion(file_{env_var_invalid}.txt) -> invalid or missing environment variable reference: {env}

}
