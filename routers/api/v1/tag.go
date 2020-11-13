package v1

import (
	"github.com/dayu26/crucian/pkg/export"

	"github.com/astaxie/beego/validation"
	"github.com/dayu26/crucian/service/tag_service"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"

	"github.com/dayu26/crucian/pkg/app"
	"github.com/dayu26/crucian/pkg/e"

	"github.com/dayu26/crucian/pkg/logging"
	"github.com/dayu26/crucian/pkg/setting"
	"github.com/dayu26/crucian/pkg/util"
)

// GetTags get tag from db
func GetTags(c *gin.Context) {
	name := c.Query("name")
	state := -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	tagService := tag_service.Tag{
		Name:     name,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}
	tags, err := tagService.GetAll()
	if err != nil {
		app.JsonError(c, e.ERROR_GET_TAGS_FAIL, nil)
		return
	}

	count, err := tagService.Count()
	if err != nil {
		app.JsonError(c, e.ERROR_COUNT_TAG_FAIL, nil)
		return
	}

	app.JsonSuccess(c, e.SUCCESS, gin.H{
		"lists": tags,
		"total": count,
	})
}

type AddTagForm struct {
	Name      string `form:"name" valid:"Required;MaxSize(100)"`
	CreatedBy string `form:"created_by" valid:"Required;MaxSize(100)"`
	State     int    `form:"state" valid:"Range(0,1)"`
}

// AddTag function
func AddTag(c *gin.Context) {
	var form AddTagForm

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		app.Json(c, httpCode, errCode, nil)
		return
	}

	tagService := tag_service.Tag{
		Name:      form.Name,
		CreatedBy: form.CreatedBy,
		State:     form.State,
	}
	exists, err := tagService.ExistByName()
	if err != nil {
		app.JsonError(c, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if exists {
		app.JsonError(c, e.ERROR_EXIST_TAG, nil)
		return
	}

	err = tagService.Add()
	if err != nil {
		app.JsonError(c, e.ERROR_ADD_TAG_FAIL, nil)
		return
	}

	app.JsonSuccess(c, e.SUCCESS, nil)
}

type EditTagForm struct {
	ID         int    `form:"id" valid:"Required;Min(1)"`
	Name       string `form:"name" valid:"Required;MaxSize(100)"`
	ModifiedBy string `form:"modified_by" valid:"Required;MaxSize(100)"`
	State      int    `form:"state" valid:"Range(0,1)"`
}

// EditTag function
func EditTag(c *gin.Context) {
	var form = EditTagForm{ID: com.StrTo(c.Param("id")).MustInt()}

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		app.Json(c, httpCode, errCode, nil)
		return
	}

	tagService := tag_service.Tag{
		ID:         form.ID,
		Name:       form.Name,
		ModifiedBy: form.ModifiedBy,
		State:      form.State,
	}

	exists, err := tagService.ExistByID()
	if err != nil {
		app.JsonError(c, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exists {
		app.JsonError(c, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = tagService.Edit()
	if err != nil {
		app.JsonError(c, e.ERROR_EDIT_TAG_FAIL, nil)
		return
	}

	app.JsonSuccess(c, e.SUCCESS, nil)
}

//DeleteTag function
func DeleteTag(c *gin.Context) {
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		app.JsonError(c, e.INVALID_PARAMS, nil)
		return
	}

	tagService := tag_service.Tag{ID: id}
	exists, err := tagService.ExistByID()
	if err != nil {
		app.JsonError(c, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exists {
		app.JsonError(c, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	if err := tagService.Delete(); err != nil {
		app.JsonError(c, e.ERROR_DELETE_TAG_FAIL, nil)
		return
	}

	app.JsonSuccess(c, e.SUCCESS, nil)
}

//ExportTag function
func ExportTag(c *gin.Context) {
	name := c.PostForm("name")
	state := -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	tagService := tag_service.Tag{
		Name:  name,
		State: state,
	}

	filename, err := tagService.Export()
	if err != nil {
		app.JsonError(c, e.ERROR_EXPORT_TAG_FAIL, nil)
		return
	}

	app.JsonSuccess(c, e.SUCCESS, gin.H{
		"export_url":      export.GetExcelFullUrl(filename),
		"export_save_url": export.GetExcelPath() + filename,
	})
}

// ImportTag function
func ImportTag(c *gin.Context) {

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		logging.Warn(err)
		app.JsonError(c, e.ERROR, nil)
		return
	}

	tagService := tag_service.Tag{}
	err = tagService.Import(file)
	if err != nil {
		logging.Warn(err)
		app.JsonError(c, e.ERROR_IMPORT_TAG_FAIL, nil)
		return
	}

	app.JsonSuccess(c, e.SUCCESS, nil)
}
