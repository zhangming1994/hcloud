package controllers

import (
	"github.com/astaxie/beego"
)

type ErrorController struct {
	beego.Controller
}

func (c *ErrorController) Error401() {
	c.Data["content"] = "page not found"
	c.Data["err"] = "401错误"
	c.TplName = "err/error404.html"
}

func (c *ErrorController) Error402() {
	c.Data["content"] = "page not found"
	c.Data["err"] = "402错误"
	c.TplName = "err/error404.html"
}

func (c *ErrorController) Error403() {
	c.Data["content"] = "page not found"
	c.Data["err"] = "403错误"
	c.TplName = "err/error404.html"
}

func (c *ErrorController) Error404() {
	c.Data["content"] = "page not found"
	c.Data["err"] = "404错误"
	c.TplName = "err/error404.html"
}

func (c *ErrorController) Error405() {
	c.Data["content"] = "page not found"
	c.Data["err"] = "405错误"
	c.TplName = "err/error404.html"
}

func (c *ErrorController) Error500() {
	c.Data["content"] = "page not found"
	c.Data["err"] = "500错误"
	c.TplName = "err/error404.html"
}

func (c *ErrorController) Error501() {
	c.Data["content"] = "page not found"
	c.Data["err"] = "501错误"
	c.TplName = "err/error404.html"
}

func (c *ErrorController) Error502() {
	c.Data["content"] = "page not found"
	c.Data["err"] = "502错误"
	c.TplName = "err/error404.html"
}

func (c *ErrorController) Error503() {
	c.Data["content"] = "page not found"
	c.Data["err"] = "503错误"
	c.TplName = "err/error404.html"
}

func (c *ErrorController) Error504() {
	c.Data["content"] = "page not found"
	c.Data["err"] = "504错误"
	c.TplName = "err/error404.html"
}
