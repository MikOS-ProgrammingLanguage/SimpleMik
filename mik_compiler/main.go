package main

import (
	"MicNewDawn/errors"
	"MicNewDawn/mik_compiler/bigO"
	executablegenerator "MicNewDawn/mik_compiler/generators/executable_generator"
	"MicNewDawn/mik_compiler/lexer"
	"MicNewDawn/mik_compiler/parser"
	"MicNewDawn/mik_compiler/preprocessor"
	"MicNewDawn/utils"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

var TARGET string
var PURE_SOURCE_CODE []string

func main() {
	// Flags used to modify the behavior of the compiler
	input_flag := flag.String("i", "", "'-i' is the input flag. Usage: -i <file_name>")
	output_flag := flag.String("o", "mik.out", "'-o' is the output flag. Your output file will have it's value. Usage: -o <file_name>")

	target_executable := flag.Bool("exec", true, "Specifies if the intended output file should be a executable. True by default")
	target_llvm_ir := flag.Bool("ll", false, "Specifies if the intended output file should be a llvm intermediate representation")
	target_obj := flag.Bool("obj", false, "Produce a linkable object file")

	emit_preprocessor := flag.Bool("emit-preprocessor", false, "Prints the preprocessed text")
	emit_lexer := flag.Bool("emit-lexer", false, "Prints the output of the lexer")

	time_steps := flag.Bool("time-steps", false, "Prints timestamps for each compilation step")

	ignore_unused := flag.Bool("ou", false, "Ignore all unused functions and variables")
	unused_warning := flag.Bool("wu", false, "Throw warning for each unused variable")

	ignore_sections_flag := flag.String("sec-ign", "", "Ignore one or more sections (separated by a comma).")
	install_flag := flag.Bool("install", false, "Install a config file and the mip package structure. Must be done before first use. Otherwise the compiler doesn't know where it is installed")
	wall_flag := flag.Bool("Wall", false, "Emit all warnings")

	big_o_flag := flag.Bool("O", false, "Emits an estimate of the big O of a function")

	// Meme flags
	fuck_you_flag := flag.Bool("fu", false, "Fucks your ass on error.")

	/*	optimization_config_file := flag.String("opt-cfg", "", "Use a specified file as a config file for optimizations")
		show_optimizations := flag.Bool("wopts", false, "Emit all the optimizations made as warnings to the terminal")
		optimize_empty_blocks_flag := flag.Bool("o-eb", false, "Optimize all empty code blocks")
		optimize_dead_code_flag := flag.Bool("o-dc", false, "Optimize all occurrences of dead code")

		// Create optimization struct
		optFlags := optimizer.OptimizationFlags{
			ShowOptimizationsMade: *show_optimizations,
			OptimizeEmptyBlocks:   *optimize_empty_blocks_flag,
			OptimizeDeadCode:      *optimize_dead_code_flag,
		}

		if *optimization_config_file != "" {
			optFlags = config.ParseOptimizationConfig(optimization_config_file)
		}
	*/
	flag.Parse()
	// validate flags
	ignored_sections := strings.Split(*ignore_sections_flag, ",")
	errors.WALL_FLAG = *wall_flag

	// validate which target should be used
	if *target_executable && *target_llvm_ir {
		TARGET = "llvm"
	} else if *target_executable {
		TARGET = "bin"
	} else if *target_obj {
		TARGET = "obj"
	}

	// Fuck you!
	if *fuck_you_flag {
		if os.Geteuid() != 0 {
			errors.ThrowError("NoMaidens. Unfortunately for you, you are maidenless 😣🤨 (use sudo)", true)
		}

		fmt.Println("Are you sure you want to proceed? [y/n]")
		var yes_no rune
		fmt.Scanf("%c", &yes_no)

		if yes_no == 'n' {
			fmt.Println("I'll only segfault :)")
			var sigsegv *string
			*sigsegv = "Hello"
		}
	}

	// check for install
	if *install_flag {
		utils.Install_config()
	} else {
		// load source code
		if *input_flag == "" {
			flag.Usage()
			os.Exit(0)
		}
		a, err := os.ReadFile(*input_flag)
		if err != nil {
			panic(err)
		}
		source_code := string(a)
		PURE_SOURCE_CODE = strings.Split(source_code, "\n")

		// start compiling
		start_time := time.Now()

		preprocessed_time_start := time.Now()
		preprocessed := preprocessor.Preprocess(&source_code, input_flag)
		preprocessed_time := time.Since(preprocessed_time_start).String()

		if *emit_preprocessor {
			errors.ThrowInfo("Preprocessed Output:", false)
		}

		lexed_time_start := time.Now()
		lexed, sections := lexer.Lex(input_flag, preprocessed, ignored_sections)
		lexed_time := time.Since(lexed_time_start).String()

		if *emit_lexer {
			errors.ThrowInfo("Lexed representation:", false)
			for _, i := range lexed {
				fmt.Printf("[%s, %d, %d, \"%s\"]\n", i.Section_name, i.Line, i.Token_type, i.Value)
			}
		}

		parser_time_start := time.Now()
		AST, knownBigO := parser.Parse(lexed, sections, PURE_SOURCE_CODE, *fuck_you_flag)
		parser_time := time.Since(parser_time_start).String()
		/*
			optimization_time_start := time.Now()
			optimizedAST := optimizer.Optimize(optFlags, 0, &AST)
			optimization_time := time.Since(optimization_time_start).String()
		*/

		if *big_o_flag {
			bigOs := bigO.EmitBigOOfAST(AST, knownBigO)
			fmt.Println(bigOs)
		}

		code_generation_time_start := time.Now()
		if TARGET == "bin" || TARGET == "llvm" || TARGET == "obj" {
			executablegenerator.GenerateExecutableOrIntermediateRepresentation(&AST, *output_flag, TARGET, *ignore_unused, *unused_warning)
		}
		code_generation_time := time.Since(code_generation_time_start).String()

		if *unused_warning {
			for _, name := range parser.NamesAndScopes.Variable_names {
				if !parser.ItemInIntMap(name, parser.NamesAndScopes.VARIABLE_REFERENCE_COUNT) {
					// is unused
					errors.UnusedVariableWarning(name)
				}
			}

			for _, name := range parser.NamesAndScopes.Function_names {
				if !parser.ItemInIntMap(name, parser.NamesAndScopes.VARIABLE_REFERENCE_COUNT) {
					// is unused
					errors.UnusedFunctionWarning(name)
				}
			}
		}

		// if timing
		if *time_steps {
			fmt.Println()
			errors.ThrowInfo(fmt.Sprintf("Preprocessing: %s", preprocessed_time), false)
			errors.ThrowInfo(fmt.Sprintf("Lexing: %s", lexed_time), false)
			errors.ThrowInfo(fmt.Sprintf("Parsing: %s", parser_time), false)
			//errors.ThrowInfo(fmt.Sprintf("Optimizing: %s", optimization_time), false)
			errors.ThrowInfo(fmt.Sprintf("Code Generation: %s\n", code_generation_time), false)
		}

		// if success
		errors.ThrowSuccess(fmt.Sprintf("Successfully compiled: %s. In %s. To %s", *input_flag, time.Since(start_time).String(), *output_flag), false)
	}
}
