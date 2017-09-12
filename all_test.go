package cwl

import (
	"fmt"
	"os"
	"testing"

	. "github.com/otiai10/mint"
)

const version = "1.0"

// Provides file object for testable official .cwl files.
func cwl(name string) *os.File {
	fpath := fmt.Sprintf("./cwl/v%[1]s/v%[1]s/%s", version, name)
	f, err := os.Open(fpath)
	if err != nil {
		panic(err)
	}
	return f
}

func TestDecode_bwa_mem_tool(t *testing.T) {
	f := cwl("bwa-mem-tool.cwl")
	root := NewCWL()
	Expect(t, root).TypeOf("*cwl.Root")
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, root.Hints).TypeOf("cwl.Hints")
	Expect(t, root.Hints[0]["class"]).ToBe("ResourceRequirement")
	Expect(t, root.Hints[0]["coresMin"]).ToBe(float64(2))

	Expect(t, len(root.Inputs)).ToBe(5)
	Expect(t, root.Inputs[0]).TypeOf("cwl.RequiredInput")
	Expect(t, root.Inputs[0].ID).ToBe("reference")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Inputs[0].Binding.Position).ToBe(2)
	Expect(t, root.Inputs[1].ID).ToBe("reads")
	Expect(t, root.Inputs[1].Types[0].Type).ToBe("array")
	Expect(t, root.Inputs[1].Types[0].Items).ToBe("File")
	Expect(t, root.Inputs[1].Binding.Position).ToBe(3)
	Expect(t, root.Inputs[2].Binding.Prefix).ToBe("-m")
	Expect(t, root.Inputs[3].Binding.Separator).ToBe(",")
	Expect(t, root.Inputs[4].Default.Class).ToBe("File")
	Expect(t, root.Inputs[4].Default.Location).ToBe("args.py")
	Expect(t, root.Outputs[0].ID).ToBe("sam")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("null")
	Expect(t, root.Outputs[0].Types[1].Type).ToBe("File")
	Expect(t, root.Outputs[0].Binding.Glob).ToBe("output.sam")
	Expect(t, root.Outputs[1].ID).ToBe("args")
	Expect(t, root.Outputs[1].Types[0].Type).ToBe("array")
	Expect(t, root.Outputs[1].Types[0].Items).ToBe("string")
}

func TestDecode_binding_test(t *testing.T) {
	f := cwl("binding-test.cwl")
	root := NewCWL()
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)

	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")

	Expect(t, root.Hints[0]["class"]).ToBe("DockerRequirement")
	Expect(t, root.Hints[0]["dockerPull"]).ToBe("python:2-slim")

	Expect(t, root.Inputs[0].ID).ToBe("reference")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Inputs[0].Binding.Position).ToBe(2)
	Expect(t, root.Inputs[1].ID).ToBe("reads")
	Expect(t, root.Inputs[1].Types[0].Type).ToBe("array")
	Expect(t, root.Inputs[1].Types[0].Items).ToBe("File")
	Expect(t, root.Inputs[1].Types[0].Binding.Prefix).ToBe("-YYY")
	Expect(t, root.Inputs[1].Binding.Position).ToBe(3)
	Expect(t, root.Inputs[1].Binding.Prefix).ToBe("-XXX")
	Expect(t, root.Inputs[2].ID).ToBe("#args.py")
	Expect(t, root.Inputs[2].Types[0].Type).ToBe("File")
	Expect(t, root.Inputs[2].Default.Class).ToBe("File")
	Expect(t, root.Inputs[2].Default.Location).ToBe("args.py")
	Expect(t, root.Inputs[2].Binding.Position).ToBe(-1)

	Expect(t, root.Outputs[0].ID).ToBe("args")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("string[]")
}

func TestDecode_tmap_tool(t *testing.T) {
	f := cwl("tmap-tool.cwl")
	root := NewCWL()
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)

	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")

	Expect(t, root.Hints[0]["class"]).ToBe("DockerRequirement")
	Expect(t, root.Hints[0]["dockerPull"]).ToBe("python:2-slim")

	Expect(t, root.Inputs[0].ID).ToBe("reads")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Inputs[1].ID).ToBe("stages")
	Expect(t, root.Inputs[1].Types[0].Type).ToBe("array")
	Expect(t, root.Inputs[1].Types[0].Items).ToBe("#Stage")
	Expect(t, root.Inputs[1].Binding.Position).ToBe(1)
	Expect(t, root.Inputs[2].ID).ToBe("#args.py")
	Expect(t, root.Inputs[2].Types[0].Type).ToBe("File")
	Expect(t, root.Inputs[2].Default.Class).ToBe("File")
	Expect(t, root.Inputs[2].Default.Location).ToBe("args.py")
	Expect(t, root.Inputs[2].Binding.Position).ToBe(-1)

	Expect(t, root.Outputs[0].ID).ToBe("sam")
	Expect(t, root.Outputs[0].Binding.Glob).ToBe("output.sam")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("null")
	Expect(t, root.Outputs[0].Types[1].Type).ToBe("File")
	Expect(t, root.Outputs[1].ID).ToBe("args")
	Expect(t, root.Outputs[1].Types[0].Type).ToBe("string[]")

	Expect(t, root.Requirements[0].Class).ToBe("SchemaDefRequirement")
	Expect(t, root.Requirements[0].Types[0].Name).ToBe("Map1")
	Expect(t, root.Requirements[0].Types[0].Type).ToBe("record")
	Expect(t, root.Requirements[0].Types[0].Fields[0].Name).ToBe("algo")
	Expect(t, root.Requirements[0].Types[0].Fields[0].Types[0].Type).ToBe("enum")
	Expect(t, root.Requirements[0].Types[0].Fields[0].Types[0].Name).ToBe("JustMap1")
	Expect(t, root.Requirements[0].Types[0].Fields[0].Types[0].Symbols[0]).ToBe("map1")
	Expect(t, root.Requirements[0].Types[0].Fields[0].Binding.Position).ToBe(0)
	Expect(t, root.Requirements[0].Types[0].Fields[1].Name).ToBe("maxSeqLen")
	Expect(t, root.Requirements[0].Types[0].Fields[1].Types[0].Type).ToBe("null")
	Expect(t, root.Requirements[0].Types[0].Fields[1].Types[1].Type).ToBe("int")
	Expect(t, root.Requirements[0].Types[0].Fields[1].Binding.Position).ToBe(2)
	Expect(t, root.Requirements[0].Types[0].Fields[1].Binding.Prefix).ToBe("--max-seq-length")
	Expect(t, root.Requirements[0].Types[0].Fields[2].Name).ToBe("minSeqLen")
	Expect(t, root.Requirements[0].Types[0].Fields[2].Types[0].Type).ToBe("null")
	Expect(t, root.Requirements[0].Types[0].Fields[2].Types[1].Type).ToBe("int")
	Expect(t, root.Requirements[0].Types[0].Fields[2].Binding.Position).ToBe(2)
	Expect(t, root.Requirements[0].Types[0].Fields[2].Binding.Prefix).ToBe("--min-seq-length")
	Expect(t, root.Requirements[0].Types[0].Fields[3].Name).ToBe("seedLength")
	Expect(t, root.Requirements[0].Types[0].Fields[3].Types[0].Type).ToBe("null")
	Expect(t, root.Requirements[0].Types[0].Fields[3].Types[1].Type).ToBe("int")
	Expect(t, root.Requirements[0].Types[0].Fields[3].Binding.Position).ToBe(2)
	Expect(t, root.Requirements[0].Types[0].Fields[3].Binding.Prefix).ToBe("--seed-length")

	Expect(t, root.Requirements[0].Types[1].Name).ToBe("Map2")
	Expect(t, root.Requirements[0].Types[1].Type).ToBe("record")
	Expect(t, root.Requirements[0].Types[1].Fields[0].Name).ToBe("algo")
	Expect(t, root.Requirements[0].Types[1].Fields[0].Types[0].Type).ToBe("enum")
	Expect(t, root.Requirements[0].Types[1].Fields[0].Types[0].Name).ToBe("JustMap2")
	Expect(t, root.Requirements[0].Types[1].Fields[0].Types[0].Symbols[0]).ToBe("map2")
	Expect(t, root.Requirements[0].Types[1].Fields[0].Binding.Position).ToBe(0)
	Expect(t, root.Requirements[0].Types[1].Fields[1].Name).ToBe("maxSeqLen")
	Expect(t, root.Requirements[0].Types[1].Fields[1].Types[0].Type).ToBe("null")
	Expect(t, root.Requirements[0].Types[1].Fields[1].Types[1].Type).ToBe("int")
	Expect(t, root.Requirements[0].Types[1].Fields[1].Binding.Position).ToBe(2)
	Expect(t, root.Requirements[0].Types[1].Fields[1].Binding.Prefix).ToBe("--max-seq-length")
	Expect(t, root.Requirements[0].Types[1].Fields[2].Name).ToBe("minSeqLen")
	Expect(t, root.Requirements[0].Types[1].Fields[2].Types[0].Type).ToBe("null")
	Expect(t, root.Requirements[0].Types[1].Fields[2].Types[1].Type).ToBe("int")
	Expect(t, root.Requirements[0].Types[1].Fields[2].Binding.Position).ToBe(2)
	Expect(t, root.Requirements[0].Types[1].Fields[2].Binding.Prefix).ToBe("--min-seq-length")
	Expect(t, root.Requirements[0].Types[1].Fields[3].Name).ToBe("maxSeedHits")
	Expect(t, root.Requirements[0].Types[1].Fields[3].Types[0].Type).ToBe("null")
	Expect(t, root.Requirements[0].Types[1].Fields[3].Types[1].Type).ToBe("int")
	Expect(t, root.Requirements[0].Types[1].Fields[3].Binding.Position).ToBe(2)
	Expect(t, root.Requirements[0].Types[1].Fields[3].Binding.Prefix).ToBe("--max-seed-hits")

	Expect(t, root.Requirements[0].Types[2].Name).ToBe("Map3")
	Expect(t, root.Requirements[0].Types[2].Type).ToBe("record")
	Expect(t, root.Requirements[0].Types[2].Fields[0].Name).ToBe("algo")
	Expect(t, root.Requirements[0].Types[2].Fields[0].Types[0].Type).ToBe("enum")
	Expect(t, root.Requirements[0].Types[2].Fields[0].Types[0].Name).ToBe("JustMap3")
	Expect(t, root.Requirements[0].Types[2].Fields[0].Types[0].Symbols[0]).ToBe("map3")
	Expect(t, root.Requirements[0].Types[2].Fields[0].Binding.Position).ToBe(0)
	Expect(t, root.Requirements[0].Types[2].Fields[1].Name).ToBe("maxSeqLen")
	Expect(t, root.Requirements[0].Types[2].Fields[1].Types[0].Type).ToBe("null")
	Expect(t, root.Requirements[0].Types[2].Fields[1].Types[1].Type).ToBe("int")
	Expect(t, root.Requirements[0].Types[2].Fields[1].Binding.Position).ToBe(2)
	Expect(t, root.Requirements[0].Types[2].Fields[1].Binding.Prefix).ToBe("--max-seq-length")
	Expect(t, root.Requirements[0].Types[2].Fields[2].Name).ToBe("minSeqLen")
	Expect(t, root.Requirements[0].Types[2].Fields[2].Types[0].Type).ToBe("null")
	Expect(t, root.Requirements[0].Types[2].Fields[2].Types[1].Type).ToBe("int")
	Expect(t, root.Requirements[0].Types[2].Fields[2].Binding.Position).ToBe(2)
	Expect(t, root.Requirements[0].Types[2].Fields[2].Binding.Prefix).ToBe("--min-seq-length")
	Expect(t, root.Requirements[0].Types[2].Fields[3].Name).ToBe("fwdSearch")
	Expect(t, root.Requirements[0].Types[2].Fields[3].Types[0].Type).ToBe("null")
	Expect(t, root.Requirements[0].Types[2].Fields[3].Types[1].Type).ToBe("boolean")
	Expect(t, root.Requirements[0].Types[2].Fields[3].Binding.Position).ToBe(2)
	Expect(t, root.Requirements[0].Types[2].Fields[3].Binding.Prefix).ToBe("--fwd-search")

	Expect(t, root.Requirements[0].Types[3].Name).ToBe("Map4")
	Expect(t, root.Requirements[0].Types[3].Type).ToBe("record")
	Expect(t, root.Requirements[0].Types[3].Fields[0].Name).ToBe("algo")
	Expect(t, root.Requirements[0].Types[3].Fields[0].Types[0].Type).ToBe("enum")
	Expect(t, root.Requirements[0].Types[3].Fields[0].Types[0].Name).ToBe("JustMap4")
	Expect(t, root.Requirements[0].Types[3].Fields[0].Types[0].Symbols[0]).ToBe("map4")
	Expect(t, root.Requirements[0].Types[3].Fields[0].Binding.Position).ToBe(0)
	Expect(t, root.Requirements[0].Types[3].Fields[1].Name).ToBe("maxSeqLen")
	Expect(t, root.Requirements[0].Types[3].Fields[1].Types[0].Type).ToBe("null")
	Expect(t, root.Requirements[0].Types[3].Fields[1].Types[1].Type).ToBe("int")
	Expect(t, root.Requirements[0].Types[3].Fields[1].Binding.Position).ToBe(2)
	Expect(t, root.Requirements[0].Types[3].Fields[1].Binding.Prefix).ToBe("--max-seq-length")
	Expect(t, root.Requirements[0].Types[3].Fields[2].Name).ToBe("minSeqLen")
	Expect(t, root.Requirements[0].Types[3].Fields[2].Types[0].Type).ToBe("null")
	Expect(t, root.Requirements[0].Types[3].Fields[2].Types[1].Type).ToBe("int")
	Expect(t, root.Requirements[0].Types[3].Fields[2].Binding.Position).ToBe(2)
	Expect(t, root.Requirements[0].Types[3].Fields[2].Binding.Prefix).ToBe("--min-seq-length")
	Expect(t, root.Requirements[0].Types[3].Fields[3].Name).ToBe("seedStep")
	Expect(t, root.Requirements[0].Types[3].Fields[3].Types[0].Type).ToBe("null")
	Expect(t, root.Requirements[0].Types[3].Fields[3].Types[1].Type).ToBe("int")
	Expect(t, root.Requirements[0].Types[3].Fields[3].Binding.Position).ToBe(2)
	Expect(t, root.Requirements[0].Types[3].Fields[3].Binding.Prefix).ToBe("--seed-step")

	Expect(t, root.Requirements[0].Types[4].Name).ToBe("Stage")
	Expect(t, root.Requirements[0].Types[4].Type).ToBe("record")
	Expect(t, root.Requirements[0].Types[4].Fields[0].Name).ToBe("stageId")
	Expect(t, root.Requirements[0].Types[4].Fields[0].Types[0].Type).ToBe("null")
	Expect(t, root.Requirements[0].Types[4].Fields[0].Binding.Position).ToBe(0)
	Expect(t, root.Requirements[0].Types[4].Fields[0].Binding.Prefix).ToBe("stage")
	Expect(t, root.Requirements[0].Types[4].Fields[0].Binding.Separate).ToBe(false)
	Expect(t, root.Requirements[0].Types[4].Fields[1].Name).ToBe("stageOption1")
	Expect(t, root.Requirements[0].Types[4].Fields[1].Types[0].Type).ToBe("null")
	Expect(t, root.Requirements[0].Types[4].Fields[1].Types[1].Type).ToBe("boolean")
	Expect(t, root.Requirements[0].Types[4].Fields[1].Binding.Position).ToBe(1)
	Expect(t, root.Requirements[0].Types[4].Fields[1].Binding.Prefix).ToBe("-n")
	Expect(t, root.Requirements[0].Types[4].Fields[2].Name).ToBe("algos")
	Expect(t, root.Requirements[0].Types[4].Fields[2].Types[0].Type).ToBe("array")
	Expect(t, root.Requirements[0].Types[4].Fields[2].Types[0].Items[0]).ToBe("#Map1")
	Expect(t, root.Requirements[0].Types[4].Fields[2].Types[0].Items[1]).ToBe("#Map2")
	Expect(t, root.Requirements[0].Types[4].Fields[2].Types[0].Items[2]).ToBe("#Map3")
	Expect(t, root.Requirements[0].Types[4].Fields[2].Types[0].Items[3]).ToBe("#Map4")
	Expect(t, root.Requirements[0].Types[4].Fields[2].Binding.Position).ToBe(2)
}

func TestDecode_cat1_testcli(t *testing.T) {
	f := cwl("cat1-testcli.cwl")
	root := NewCWL()
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)

	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, root.Doc).ToBe("Print the contents of a file to stdout using 'cat' running in a docker container.")

	Expect(t, root.Hints[0]["class"]).ToBe("DockerRequirement")
	Expect(t, root.Hints[0]["dockerPull"]).ToBe("python:2-slim")

	Expect(t, root.Inputs[0].ID).ToBe("file1")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Inputs[0].Binding.Position).ToBe(1)
	Expect(t, root.Inputs[1].ID).ToBe("numbering")
	Expect(t, root.Inputs[1].Types[0].Type).ToBe("null")
	Expect(t, root.Inputs[1].Types[1].Type).ToBe("boolean")
	Expect(t, root.Inputs[2].ID).ToBe("args.py")
	Expect(t, root.Inputs[2].Types[0].Type).ToBe("File")
	Expect(t, root.Inputs[2].Default.Class).ToBe("File")
	Expect(t, root.Inputs[2].Default.Location).ToBe("args.py")
	Expect(t, root.Inputs[2].Binding.Position).ToBe(-1)

	Expect(t, root.Outputs[0].ID).ToBe("args")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("string[]")

	Expect(t, root.BaseCommands[0]).ToBe("python")
	Expect(t, root.Arguments[0]).ToBe("cat")
}

func TestDecode_template_tool(t *testing.T) {
	f := cwl("template-tool.cwl")
	root := NewCWL()
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)

	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")

	Expect(t, root.Requirements[0].Class).ToBe("InlineJavascriptRequirement")
	Expect(t, root.Requirements[0].ExpressionLib[0].Include).ToBe("underscore.js")
	Expect(t, root.Requirements[0].ExpressionLib[1].Execute).ToBe("var t = function(s) { return _.template(s)({'inputs': inputs}); };")

	Expect(t, root.Requirements[1].Class).ToBe("InitialWorkDirRequirement")
	Expect(t, root.Requirements[1].Listing[0].Name).ToBe("foo.txt")
	Expect(t, root.Requirements[1].Listing[0].Entry).ToBe(`$(t("The file is <%= inputs.file1.path.split('/').slice(-1)[0] %>\n"))`)

	Expect(t, root.Hints[0]["class"]).ToBe("DockerRequirement")
	Expect(t, root.Hints[0]["dockerPull"]).ToBe("debian:8")

	Expect(t, root.Inputs[0].ID).ToBe("file1")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File")

	Expect(t, root.Outputs[0].ID).ToBe("foo")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Outputs[0].Binding.Glob).ToBe("foo.txt")

	Expect(t, root.BaseCommands[0]).ToBe("cat")
	Expect(t, root.BaseCommands[1]).ToBe("foo.txt")
}

func TestDecode_count_lines1_wf(t *testing.T) {
	f := cwl("count-lines1-wf.cwl")
	root := NewCWL()
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)

	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("Workflow")

	Expect(t, root.Inputs[0].ID).ToBe("file1")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Outputs[0].ID).ToBe("count_output")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("int")
	Expect(t, root.Outputs[0].Source).ToBe("step2/output")

	Expect(t, root.Steps[0].ID).ToBe("step1")
	Expect(t, root.Steps[0].Run).ToBe("wc-tool.cwl")
	Expect(t, root.Steps[0].In[0].Name).ToBe("file1")
	Expect(t, root.Steps[0].In[0].Location).ToBe("file1")
	Expect(t, root.Steps[0].Out[0].Name).ToBe("output")
	Expect(t, root.Steps[1].ID).ToBe("step2")
	Expect(t, root.Steps[1].Run).ToBe("parseInt-tool.cwl")
	Expect(t, root.Steps[1].In[0].Name).ToBe("file1")
	Expect(t, root.Steps[1].In[0].Location).ToBe("step1/output")
	Expect(t, root.Steps[1].Out[0].Name).ToBe("output")
}

func TestDecode_cat3_nodocker(t *testing.T) {
	f := cwl("cat3-nodocker.cwl")
	root := NewCWL()
	Expect(t, root).TypeOf("*cwl.Root")
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Doc).ToBe("Print the contents of a file to stdout using 'cat'.")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, len(root.BaseCommands)).ToBe(1)
	Expect(t, root.BaseCommands[0]).ToBe("cat")
	Expect(t, root.Stdout).ToBe("output.txt")
	Expect(t, len(root.Inputs)).ToBe(1)
	Expect(t, root.Inputs[0].ID).ToBe("file1")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Inputs[0].Label).ToBe("Input File")
	Expect(t, root.Inputs[0].Doc).ToBe("The file that will be copied using 'cat'")
	Expect(t, root.Inputs[0].Binding.Position).ToBe(1)
}

func TestDecode_cat3_tool_mediumcut(t *testing.T) {
	f := cwl("cat3-tool-mediumcut.cwl")
	root := NewCWL()
	Expect(t, root).TypeOf("*cwl.Root")
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Doc).ToBe("Print the contents of a file to stdout using 'cat' running in a docker container.")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, len(root.BaseCommands)).ToBe(1)
	Expect(t, root.BaseCommands[0]).ToBe("cat")
	Expect(t, root.Stdout).ToBe("cat-out")
	Expect(t, root.Hints).TypeOf("cwl.Hints")
	Expect(t, root.Hints[0]["class"]).ToBe("DockerRequirement")
	Expect(t, root.Hints[0]["dockerPull"]).ToBe("debian:wheezy")
	Expect(t, len(root.Inputs)).ToBe(1)
	Expect(t, root.Inputs[0].ID).ToBe("file1")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Inputs[0].Label).ToBe("Input File")
	Expect(t, root.Inputs[0].Doc).ToBe("The file that will be copied using 'cat'")
	Expect(t, root.Inputs[0].Binding.Position).ToBe(1)
}

func TestDecode_cat3_tool_shortcut(t *testing.T) {
	f := cwl("cat3-tool-shortcut.cwl")
	root := NewCWL()
	Expect(t, root).TypeOf("*cwl.Root")
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Doc).ToBe("Print the contents of a file to stdout using 'cat' running in a docker container.")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, len(root.BaseCommands)).ToBe(1)
	Expect(t, root.BaseCommands[0]).ToBe("cat")
	Expect(t, root.Hints).TypeOf("cwl.Hints")
	Expect(t, root.Hints[0]["class"]).ToBe("DockerRequirement")
	Expect(t, root.Hints[0]["dockerPull"]).ToBe("debian:wheezy")
	Expect(t, len(root.Inputs)).ToBe(1)
	Expect(t, root.Inputs[0].ID).ToBe("file1")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Inputs[0].Label).ToBe("Input File")
	Expect(t, root.Inputs[0].Doc).ToBe("The file that will be copied using 'cat'")
	Expect(t, root.Inputs[0].Binding.Position).ToBe(1)
}

func TestDecode_cat3_tool(t *testing.T) {
	f := cwl("cat3-tool.cwl")
	root := NewCWL()
	Expect(t, root).TypeOf("*cwl.Root")
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Doc).ToBe("Print the contents of a file to stdout using 'cat' running in a docker container.")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, len(root.BaseCommands)).ToBe(1)
	Expect(t, root.BaseCommands[0]).ToBe("cat")
	Expect(t, root.Stdout).ToBe("output.txt")
	Expect(t, root.Hints).TypeOf("cwl.Hints")
	Expect(t, root.Hints[0]["class"]).ToBe("DockerRequirement")
	Expect(t, root.Hints[0]["dockerPull"]).ToBe("debian:wheezy")
	Expect(t, len(root.Inputs)).ToBe(1)
	Expect(t, root.Inputs[0].ID).ToBe("file1")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Inputs[0].Label).ToBe("Input File")
	Expect(t, root.Inputs[0].Doc).ToBe("The file that will be copied using 'cat'")
	Expect(t, root.Inputs[0].Binding.Position).ToBe(1)
}

func TestDecode_env_tool1(t *testing.T) {
	f := cwl("env-tool1.cwl")
	root := NewCWL()
	Expect(t, root).TypeOf("*cwl.Root")
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, len(root.BaseCommands)).ToBe(3)
	Expect(t, root.BaseCommands[0]).ToBe("/bin/bash")
	Expect(t, root.BaseCommands[1]).ToBe("-c")
	Expect(t, root.BaseCommands[2]).ToBe("echo $TEST_ENV")
	Expect(t, len(root.Inputs)).ToBe(1)
	// TODO ignore "in: string'
	Expect(t, len(root.Outputs)).ToBe(1)
	Expect(t, root.Outputs[0].ID).ToBe("out")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Outputs[0].Binding.Glob).ToBe("out")
}

func TestDecode_default_path(t *testing.T) {
	f := cwl("default_path.cwl")
	root := NewCWL()
	Expect(t, root).TypeOf("*cwl.Root")
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, len(root.Inputs)).ToBe(1)
	Expect(t, root.Inputs[0].ID).ToBe("file1")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File")
	// TODO support default: section
	// TODO support outputs: []
	Expect(t, len(root.Arguments)).ToBe(2)
	Expect(t, root.Arguments[0]).ToBe("cat")
	Expect(t, root.Arguments[1]).ToBe("$(inputs.file1.path)")
}

func TestDecode_cat4_tool(t *testing.T) {
	f := cwl("cat4-tool.cwl")
	root := NewCWL()
	Expect(t, root).TypeOf("*cwl.Root")
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, len(root.Inputs)).ToBe(1)
	Expect(t, root.Inputs[0].ID).ToBe("file1")
	Expect(t, len(root.Outputs)).ToBe(1)
	Expect(t, root.Outputs[0].ID).ToBe("output_txt")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Outputs[0].Binding.Glob).ToBe("output.txt")
	Expect(t, len(root.BaseCommands)).ToBe(1)
	Expect(t, root.BaseCommands[0]).ToBe("cat")
	Expect(t, root.Stdout).ToBe("output.txt")
	Expect(t, root.Stdin).ToBe("$(inputs.file1.path)")
}

func TestDecode_cat5_tool(t *testing.T) {
	f := cwl("cat5-tool.cwl")
	root := NewCWL()
	Expect(t, root).TypeOf("*cwl.Root")
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, root.Doc).ToBe("Print the contents of a file to stdout using 'cat' running in a docker container.")
	Expect(t, len(root.Hints)).ToBe(2)
	Expect(t, root.Hints).TypeOf("cwl.Hints")
	Expect(t, root.Hints[0]["class"]).ToBe("DockerRequirement")
	Expect(t, root.Hints[0]["dockerPull"]).ToBe("debian:wheezy")
	Expect(t, root.Hints[1]["class"]).ToBe("ex:BlibberBlubberFakeRequirement")
	Expect(t, root.Hints[1]["fakeField"]).ToBe("fraggleFroogle")
	Expect(t, len(root.Inputs)).ToBe(1)
	Expect(t, root.Inputs[0].ID).ToBe("file1")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Inputs[0].Label).ToBe("Input File")
	Expect(t, root.Inputs[0].Doc).ToBe("The file that will be copied using 'cat'")
	Expect(t, root.Inputs[0].Binding.Position).ToBe(1)
	Expect(t, len(root.Outputs)).ToBe(1)
	Expect(t, root.Outputs[0].ID).ToBe("output_file")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Outputs[0].Binding.Glob).ToBe("output.txt")
	Expect(t, len(root.BaseCommands)).ToBe(1)
	Expect(t, root.BaseCommands[0]).ToBe("cat")
	Expect(t, root.Stdout).ToBe("output.txt")
	// $namespaces
	Expect(t, len(root.Namespaces)).ToBe(1)
	Expect(t, root.Namespaces[0]["ex"]).ToBe("http://example.com/")
}

func TestDecode_metadata(t *testing.T) {
	f := cwl("metadata.cwl")
	root := NewCWL()
	Expect(t, root).TypeOf("*cwl.Root")
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, root.Doc).ToBe("Print the contents of a file to stdout using 'cat' running in a docker container.")
	Expect(t, len(root.Hints)).ToBe(1)
	Expect(t, root.Hints).TypeOf("cwl.Hints")
	Expect(t, root.Hints[0]["class"]).ToBe("DockerRequirement")
	Expect(t, root.Hints[0]["dockerPull"]).ToBe("debian:wheezy")
	Expect(t, len(root.Inputs)).ToBe(2)
	Expect(t, root.Inputs[0].ID).ToBe("file1")
	Expect(t, root.Inputs[0].Types[0].Type).ToBe("File")
	Expect(t, root.Inputs[0].Binding.Position).ToBe(1)
	Expect(t, len(root.Outputs)).ToBe(0)
	Expect(t, len(root.BaseCommands)).ToBe(1)
	Expect(t, root.BaseCommands[0]).ToBe("cat")
	// $namespaces
	Expect(t, len(root.Namespaces)).ToBe(2)
	Expect(t, root.Namespaces[0]["dct"]).ToBe("http://purl.org/dc/terms/")
	Expect(t, root.Namespaces[1]["foaf"]).ToBe("http://xmlns.com/foaf/0.1/")
	// $namespaces
	Expect(t, len(root.Schemas)).ToBe(2)
	Expect(t, root.Schemas[0]).ToBe("foaf.rdf")
	Expect(t, root.Schemas[1]).ToBe("dcterms.rdf")
}