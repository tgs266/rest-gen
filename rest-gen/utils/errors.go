package utils

import "github.com/dave/jennifer/jen"

func WriteErrorCheck(inputName string, msg string) jen.Code {
	return jen.If(jen.Id(inputName).Op("!=").Nil()).Block(
		jen.Panic(jen.Lit(msg)),
	)
}
