package service

import (
	. "easy-workflow/pkg/dao"
	. "easy-workflow/pkg/workflow/engine"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//这是一个增效用的方法:
//执行SQL，将结果集填充到指定struct
//struct以json方式返回
func ExecSQLThenReturnResponse(c *gin.Context, SQL string, Result interface{}, Params ...interface{}) {
	if result, err := ExecSQL(SQL, Result, Params...); err == nil {
		c.JSON(200, result)
	} else {
		c.JSON(400, err.Error()) //http code:400 错误请求 — 请求中有语法问题，或不能满足请求。
	}
}

/*
swagger注解描述 https://github.com/swaggo/swag/blob/master/README_zh-CN.md
@Summary	摘要
@Produce	API 可以产生的 MIME 类型的列表，MIME 类型你可以简单的理解为响应类型，例如：json、xml、html 等等,详细如下：
        ---Alias-------------------------MIME Type------------------------------
           json	                         application/json
           x-www-form-urlencoded	     application/x-www-form-urlencoded
           xml	                         text/xml
           plain	                     text/plain
           html	                         text/html
           mpfd	                         multipart/form-data
           json-api	                     application/vnd.api+json
           json-stream	                 application/x-json-stream
           octet-stream	                 application/octet-stream
           png	                         image/png
           jpeg	                         image/jpeg
           gif	                         image/gif
@Param	参数格式，从左到右分别为：参数名、入参类型、数据类型、是否必填、注释、example(示例)
        -入参类型有以下几种：path query header cookie  body formData
        -数据类型有 string int uint uint32 uint64 float32 bool 以及用户自定义类型(struct)
@Success	响应成功，从左到右分别为：状态码、参数类型、数据类型、注释
@Failure	响应失败，从左到右分别为：状态码、参数类型、数据类型、注释
@Router	路由，从左到右分别为：路由地址，HTTP 方法
*/

// @Summary      流程生成/升级
// @Description
// @Tags         流程定义
// @Produce      json
// @Param        ProcessName  formData string  true  "流程名称" example(员工请假)
// @Param        Resource  formData string  true  "流程定义资源(json)" example(json字符串)
// @Param        CreateUserID  formData string  true  "创建者ID" example(0001)
// @Param        Source  formData string  true  "来源" example(办公系统)
// @Success      200  {object}  int 流程ID
// @Failure      400  {object}  string 报错信息
// @Router       /process/def/save [post]
func ProcDef_Save(c *gin.Context) {
	ProcessName := c.PostForm("ProcessName")
	Resource := c.PostForm("Resource")
	CreateUserID := c.PostForm("CreateUserID")
	Source := c.PostForm("Source")

	if ProcID, err := ProcessSave(ProcessName, Resource, CreateUserID, Source); err == nil {
		c.JSON(http.StatusOK, ProcID)
	} else {
		c.JSON(400, err)
	}
}

// @Summary      获取特定source下所有流程
// @Description  引擎可能被多个系统、组件等使用，source表示从哪个来源创建的流程
// @Tags         流程定义
// @Produce      json
// @Param        source  query string  true  "来源" example(办公系统)
// @Success      200  {object}  []model.ProcessDefine 流程定义列表
// @Failure      400  {object}  string 报错信息
// @Router       /process/def/list [get]
func ProcDef_ListBySource(c *gin.Context) {
	source := c.Query("source")
	if procDef, err := GetProcessList(source); err == nil {
		c.JSON(200, procDef)
	} else {
		c.JSON(400, err)
	}
}

// @Summary      获取流程定义
// @Description  返回的是Node数组，流程是由N个Node组成的
// @Tags         流程定义
// @Produce      json
// @Param        id  query string  true  "流程ID" example(1)
// @Success      200  {object}  []model.Node "Node数组"
// @Failure      400  {object}  string 报错信息
// @Router       /process/def/get [get]
func ProcDef_GetProcDefByID(c *gin.Context) {
	id := c.Query("id")
	id_int, _ := strconv.Atoi(id)
	if nodes, err := GetProcessDefine(id_int); err == nil {
		c.JSON(200, nodes)
	} else {
		c.JSON(400, err)
	}
}

// @Summary      开始流程
// @Description  返回流程实例ID
// @Tags         流程实例
// @Produce      json
// @Param        ProcessID  formData string  true  "流程ID" example(1)
// @Param        BusinessID  formData string  true  "业务员ID" example("订单001")
// @Param        Comment  formData string  false  "评论意见" example("家中有事请假三天,请领导批准")
// @Param        VariablesJson  formData string  false  "变量(Json)" example("{"User":"001"}")
// @Success      200  {object}  int 流程实例ID
// @Failure      400  {object}  string 报错信息
// @Router       /process/inst/start [post]
func ProcInst_Start(c *gin.Context) {
	ProcessID, _ := strconv.Atoi(c.PostForm("ProcessID"))
	BusinessID := c.PostForm("BusinessID")
	Comment := c.PostForm("Comment")
	VariablesJson := make(map[string]string)
	json.Unmarshal([]byte(c.PostForm("Variables")), VariablesJson)

	if id, err := InstanceStart(ProcessID, BusinessID, Comment, VariablesJson); err == nil {
		c.JSON(200, id)
	} else {
		c.JSON(400, err)
	}
}

// @Summary      任务通过
// @Description
// @Tags         任务
// @Produce      json
// @Param        TaskID  formData string  true  "任务ID" example(1)
// @Param        Comment  formData string  false  "评论意见" example("同意请假")
// @Param        VariablesJson  formData string  false  "变量(Json)" example("{"User":"001"}")
// @Success      200  {object}  string "ok"
// @Failure      400  {object}  string 报错信息
// @Router       /process/task/pass [post]
func Task_Pass(c *gin.Context) {
	TaskID, _ := strconv.Atoi(c.PostForm("TaskID"))
	Comment := c.PostForm("Comment")
	VariableJson := c.PostForm("VariableJson")

	if err := TaskPass(TaskID, Comment, VariableJson); err == nil {
		c.JSON(200, "ok")
	} else {
		c.JSON(400, err)
	}
}

// @Summary      任务驳回
// @Description
// @Tags         任务
// @Produce      json
// @Param        TaskID  formData string  true  "任务ID" example(1)
// @Param        Comment  formData string  false  "评论意见" example("不同意")
// @Param        VariablesJson  formData string  false  "变量(Json)" example("{"User":"001"}")
// @Success      200  {object}  string "ok"
// @Failure      400  {object}  string 报错信息
// @Router       /process/task/reject [post]
func Task_Reject(c *gin.Context) {
	TaskID, _ := strconv.Atoi(c.PostForm("TaskID"))
	Comment := c.PostForm("Comment")
	VariableJson := c.PostForm("VariableJson")

	if err := TaskReject(TaskID, Comment, VariableJson); err == nil {
		c.JSON(200, "ok")
	} else {
		c.JSON(400, err)
	}
}

// @Summary      获取待办任务
// @Description  返回的是任务数组
// @Tags         任务
// @Produce      json
// @Param        userid  query string  true  "用户ID" example("U001")
// @Success      200  {object}  []model.Task 任务数组
// @Failure      400  {object}  string 报错信息
// @Router       /process/task/todo [get]
func Task_ToDoList(c *gin.Context) {
	UserID := c.Query("userid")
	if tasks, err := GetTaskToDoList(UserID); err == nil {
		c.JSON(200, tasks)
	} else {
		c.JSON(400, err)
	}
}


// @Summary      获取已办任务
// @Description  返回的是任务数组
// @Tags         任务
// @Produce      json
// @Param        userid  query string  true  "用户ID" example("U001")
// @Success      200  {object}  []model.Task 任务数组
// @Failure      400  {object}  string 报错信息
// @Router       /process/task/todo [get]
func Task_FinishedList(c *gin.Context) {
	UserID := c.Query("userid")
	if tasks, err := GetTaskFinishedList(UserID); err == nil {
		c.JSON(200, tasks)
	} else {
		c.JSON(400, err)
	}
}