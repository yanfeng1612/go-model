package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/yanfeng1612/go-model/utils"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Config struct {
	BackTheme        string        `json:"backTheme"`
	FrontTheme       string        `json:"frontTheme"`
	ProjectName      string        `json:"projectName"`
	ProjectAdminName string        `json:"projectAdminName"`
	PackageName      string        `json:"packageName"`
	ExportPath       string        `json:"exportPath"`
	DataBaseName     string        `json:"dataBaseName"`
	DatabaseIp       string        `json:"databaseIp"`
	DatabasePort     int           `json:"databasePort"`
	DatabaseUsername string        `json:"databaseUsername"`
	DatabasePassword string        `json:"databasePassword"`
	TableConfigList  []TableConfig `json:"tableConfigList"`
	FirstInit        bool          `json:"firstInit"`
	Auto             bool          `json:"auto"`
}

type TableConfig struct {
	TableName       string   `json:"tableName"`
	ShardingKey     string   `json:"shardingKey"`
	DoubleWrite     int      `json:"doubleWrite"`
	SingleQueryList []string `json:"singleQueryList"`
	InListQuery     []string `json:"inListQuery"`
	UpdateFieldList []string `json:"updateFieldList"`
}

type Theme struct {
	ThemeName string
	templates []string
}

const (
	ConfigPath             = "D://go//src//github.com//yanfeng1612//go-model//go-model.json" // 配置文件路径
	TemplateRootPath       = "D://go//src//github.com//yanfeng1612//go-model//template"      // 模板根路径
	TemplateJavaRootPath   = TemplateRootPath + "//java"                                     // java模板根路径
	TemplateVueRootPath    = TemplateRootPath + "//vue"                                      // java模板根路径
	PomPath                = TemplateJavaRootPath + "/utils/pom.vm"                          // pom路径
	BasicResultPath        = TemplateJavaRootPath + "/utils/BasicResult.vm"                  // BasicResultPath路径
	GenericResultPath      = TemplateJavaRootPath + "/utils/GenericResult.vm"                // GenericResult路径
	CodeEnumPath           = TemplateJavaRootPath + "/utils/CodeEnum.vm"                     // CodeEnumPath路径
	QueryPath              = TemplateJavaRootPath + "/utils/Query.vm"                        // Query路径
	PageQueryPath          = TemplateJavaRootPath + "/utils/PageQuery.vm"                    // PageQuery路径
	IdPageQueryPath        = TemplateJavaRootPath + "/utils/IdPageQuery.vm"                  // IdPageQuery路径
	PagenationPath         = TemplateJavaRootPath + "/utils/Pagenation.vm"                   // Pagenation路径
	PageQueryWrapperPath   = TemplateJavaRootPath + "/utils/PageQueryWrapper.vm"             // PageQueryWrapper路径
	ListResultPath         = TemplateJavaRootPath + "/utils/ListResult.vm"                   // ListResult路径
	PageListResultPath     = TemplateJavaRootPath + "/utils/PageListResult.vm"               // PageListResult路径
	APIEmRequestStatusPath = TemplateJavaRootPath + "/utils/APIEmRequestStatus.vm"           // APIEmRequestStatus路径
	APIMsgCodePath         = TemplateJavaRootPath + "/utils/APIMsgCode.vm"                   // APIMsgCode路径
	ResponsePath           = TemplateJavaRootPath + "/utils/Response.vm"                     // Response路径
	ResponseTemplatePath   = TemplateJavaRootPath + "/utils/ResponseTemplate.vm"             // ResponseTemplate路径
	CodeConverterPath      = TemplateJavaRootPath + "/utils/CodeConverter.vm"                // CodeConverter路径
	BootstrapPath          = TemplateJavaRootPath + "/utils/Bootstrap.vm"                    // Bootstrap路径
	ApplicationConfigPath  = TemplateJavaRootPath + "/utils/application.vm"                  // application路径

	PojoPath        = TemplateJavaRootPath + "/pojo.vm"        // pojo路径
	PojoQueryPath   = TemplateJavaRootPath + "/pojoQuery.vm"   // pojoQuery路径
	ControllerPath  = TemplateJavaRootPath + "/controller.vm"  // controller路径
	ServicePath     = TemplateJavaRootPath + "/service.vm"     // service路径
	ServiceImplPath = TemplateJavaRootPath + "/serviceImpl.vm" // serviceImpl路径
	ManagerPath     = TemplateJavaRootPath + "/manager.vm"     // manager路径
	ManagerImplPath = TemplateJavaRootPath + "/managerImpl.vm" // managerImpl路径
	MapperPath      = TemplateJavaRootPath + "/mapper.vm"      // mapper路径
	MapperXmlPath   = TemplateJavaRootPath + "/mapperXml.vm"   // mapperXml路径
	vueListPath     = TemplateVueRootPath + "/list.vm"         // viewList路径
	vueEditPath     = TemplateVueRootPath + "/edit.vm"         // viewEdit路径
	vueDetailPath   = TemplateVueRootPath + "/detail.vm"       // viewDetail路径
	vueIndexPath    = TemplateVueRootPath + "/index.vm"        // viewIndex路径
	vueSidebarPath  = TemplateVueRootPath + "/sidebar.vm"      // viewSidebar路径
)

func main() {
	args := os.Args
	configPath := ConfigPath
	if len(args) > 1 {
		configPath = args[1]
	} else {

	}
	jsonTxt, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	config := new(Config)
	err = json.Unmarshal(jsonTxt, &config)
	if err != nil {
		panic(err)
	}
	jsonString, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonString))

	prepare(*config)

	schema := utils.ParseSqlFromMySqlSchema(config.DatabaseIp, config.DatabasePort, config.DatabaseUsername, config.DatabasePassword, config.DataBaseName)

	inListMap := make(map[string]int)
	singleQueryMap := make(map[string]int)
	updateFieldMap := make(map[string]int)
	for _, t := range config.TableConfigList {
		for _, c := range t.InListQuery {
			inListMap[t.TableName+"_"+c] = 1
		}
		for _, c := range t.SingleQueryList {
			singleQueryMap[t.TableName+"_"+c] = 1
		}
		for _, c := range t.UpdateFieldList {
			updateFieldMap[t.TableName+"_"+c] = 1
		}
	}

	for _, t := range schema.Tables {
		for _, c := range t.Columns {
			if _, ok := inListMap[t.Name+"_"+c.Name]; ok {
				c.InListQuery = 1
			} else {
				c.InListQuery = 0
			}
			if _, ok := singleQueryMap[t.Name+"_"+c.Name]; ok {
				c.SingleQuery = 1
			} else {
				c.SingleQuery = 0
			}
			if _, ok := updateFieldMap[t.Name+"_"+c.Name]; ok {
				c.WhereUpdateField = 1
			} else {
				c.WhereUpdateField = 0
			}
		}
	}

	buildUtils(*config)
	buildDynamic(*config, schema)
}

// 初始化 1.初始化文件目录
func prepare(config Config) {
	if config.BackTheme != "None" {
		path := config.ExportPath + "/" + config.ProjectName + "/src/main/java/" + strings.ReplaceAll(config.PackageName, ".", "/")
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			panic("创建项目文件夹失败!")
		}
		err = os.MkdirAll(path+"/utils/", os.ModePerm)
		if err != nil {
			panic("创建utils文件夹失败!" + err.Error())
		}

		pathSlice := []string{path + "/domain/auto/", path + "/domain/custom/",
			path + "/controller/auto/", path + "/controller/custom/",
			path + "/service/auto/", path + "/service/custom/",
			path + "/service/impl/auto/", path + "/service/impl/custom/",
			path + "/manager/auto/", path + "/manager/custom/",
			path + "/manager/impl/auto/", path + "/manager/impl/custom/",
			path + "/dao/auto/", path + "/dao/custom/",
			config.ExportPath + "/" + config.ProjectName + "/src/main/resources/mapper/auto/", config.ExportPath + "/" + config.ProjectName + "/src/main/resources/mapper/custom/",
		}

		if config.Auto {
			for _, p := range pathSlice {
				err = os.MkdirAll(p, os.ModePerm)
				if err != nil {
					panic("创建文件夹失败!" + err.Error())
				}
			}
		} else {
			err = os.MkdirAll(path+"/domain/", os.ModePerm)
			if err != nil {
				panic("创建pojo文件夹失败!" + err.Error())
			}
			err = os.MkdirAll(path+"/controller/", os.ModePerm)
			if err != nil {
				panic("创建controller失败")
			}
			err = os.MkdirAll(path+"/service/", os.ModePerm)
			if err != nil {
				panic("创建service文件夹失败!" + err.Error())
			}
			err = os.MkdirAll(path+"/service/impl/", os.ModePerm)
			if err != nil {
				panic("创建serviceImpl文件夹失败!" + err.Error())
			}
			err = os.MkdirAll(path+"/manager/", os.ModePerm)
			if err != nil {
				panic("创建manager文件夹失败!" + err.Error())
			}
			err = os.MkdirAll(path+"/manager/impl/", os.ModePerm)
			if err != nil {
				panic("创建managerImpl文件夹失败!" + err.Error())
			}
			err = os.MkdirAll(path+"/dao/", os.ModePerm)
			if err != nil {
				panic("创建dao文件夹失败!" + err.Error())
			}
			err = os.MkdirAll(config.ExportPath+"/"+config.ProjectName+"/src/main/resources/mapper/", os.ModePerm)
			if err != nil {
				panic("创建mapperXml文件夹失败!" + err.Error())
			}
		}
	}

	if config.FrontTheme != "None" {
		// 前端vue框架
		err := os.MkdirAll(config.ExportPath+"/"+config.ProjectAdminName, os.ModePerm)
		if err != nil {
			panic("创建vue-admin失败")
		}
		err = os.MkdirAll(config.ExportPath+"/"+config.ProjectAdminName+"/src", os.ModePerm)
		if err != nil {
			panic("创建vue src失败")
		}
		err = utils.CopyDir(TemplateVueRootPath+"/static", config.ExportPath+config.ProjectAdminName)
		if err != nil {
			panic("创建admin静态资源文件失败")
		}
	}
}

func buildUtils(config Config) {
	m := buildCommonMap(config)

	if config.BackTheme == "None" {
		return
	}

	prefix := config.ExportPath + "/" + config.ProjectName + "/src/main/java/" + strings.ReplaceAll(config.PackageName, ".", "/")
	if config.FirstInit {
		run(PomPath, config.ExportPath+"/"+config.ProjectName+"/pom.xml", m)
		run(BootstrapPath, prefix+"/Bootstrap.java", m)
		run(ApplicationConfigPath, config.ExportPath+"/"+config.ProjectName+"/src/main/resources/application.yml", m)
	}
	run(BasicResultPath, prefix+"/utils/BasicResult.java", m)
	run(CodeEnumPath, prefix+"/utils/CodeEnum.java", m)
	run(QueryPath, prefix+"/utils/Query.java", m)
	run(GenericResultPath, prefix+"/utils/GenericResult.java", m)
	run(PagenationPath, prefix+"/utils/Pagenation.java", m)
	run(PageQueryPath, prefix+"/utils/PageQuery.java", m)
	run(PageQueryWrapperPath, prefix+"/utils/PageQueryWrapper.java", m)
	run(ListResultPath, prefix+"/utils/ListResult.java", m)
	run(PageListResultPath, prefix+"/utils/PageListResult.java", m)
	run(IdPageQueryPath, prefix+"/utils/IdPageQuery.java", m)
	run(APIEmRequestStatusPath, prefix+"/utils/APIEmRequestStatus.java", m)
	run(APIMsgCodePath, prefix+"/utils/APIMsgCode.java", m)
	run(ResponsePath, prefix+"/utils/Response.java", m)
	run(ResponseTemplatePath, prefix+"/utils/ResponseTemplate.java", m)
	run(CodeConverterPath, prefix+"/utils/CodeConverter.java", m)
}

func run(templatePath, outputPath string, m map[string]interface{}) {
	t := buildCommon(templatePath)
	tmpFilePath := outputPath + "Tmp"
	tmpFile, err := os.OpenFile(tmpFilePath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 066)
	err = t.Execute(tmpFile, m)
	if err != nil {
		panic(err)
	}
	err = tmpFile.Close()
	if err != nil {
		panic(err)
	}

	tmp, err := os.OpenFile(tmpFilePath, os.O_RDONLY, 066)
	file, err := os.OpenFile(outputPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 066)
	defer file.Close()
	reader := bufio.NewReader(tmp)
	preLine := ""
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		newLine := strings.ReplaceAll(string(line), "&lt;", "<")
		if strings.TrimSpace(newLine) == "" && preLine == "" {
			continue
		}
		_, err = file.WriteString(newLine + "\r\n")
		if err != nil {
			panic(err)
		}
		preLine = strings.TrimSpace(newLine)
	}
	err = tmp.Close()
	if err != nil {
		panic(err)
	}
	err = os.Remove(tmpFilePath)
	if err != nil {
		panic(err)
	}
}

func buildCommon(templatePath string) *template.Template {
	textByte, err := ioutil.ReadFile(templatePath)
	if err != nil {
		panic(err)
	}
	text := string(textByte)
	t := template.New("")
	if t, err = t.Parse(text); err != nil {
		panic(err)
	}
	return t
}

func buildCommonMap(config Config) map[string]interface{} {
	m := make(map[string]interface{})
	m["author"] = "model-driven"
	now := time.Now()
	m["date"] = now.Format("2006-01-02")
	m["config"] = config
	return m
}

func buildDynamic(config Config, schema *utils.Schema) {
	var pathPrefix = config.ExportPath + "/" + config.ProjectName + "/src/main/java/" + strings.ReplaceAll(config.PackageName, ".", "/")
	for _, table := range schema.Tables {
		m := buildCommonMap(config)
		m["backendObject"] = table
		m["columnSize"] = len(table.Columns)

		if config.BackTheme != "None" {
			if config.Auto {
				run(PojoPath, pathPrefix+"/domain/auto/"+table.JavaBeanName+".java", m)
				run(PojoQueryPath, pathPrefix+"/domain/auto/"+table.JavaBeanName+"Query.java", m)
				run(ControllerPath, pathPrefix+"/controller/auto/"+table.JavaBeanName+"Controller.java", m)
				run(ServicePath, pathPrefix+"/service/auto/"+table.JavaBeanName+"Service.java", m)
				run(ServiceImplPath, pathPrefix+"/service/impl/auto/"+table.JavaBeanName+"ServiceImpl.java", m)
				run(ManagerPath, pathPrefix+"/manager/auto/"+table.JavaBeanName+"Manager.java", m)
				run(ManagerImplPath, pathPrefix+"/manager/impl/auto/"+table.JavaBeanName+"ManagerImpl.java", m)
				run(MapperPath, pathPrefix+"/dao/auto/"+table.JavaBeanName+"Mapper.java", m)
				run(MapperXmlPath, config.ExportPath+"/"+config.ProjectName+"/src/main/resources/mapper/auto/"+table.JavaBeanName+"Mapper.xml", m)
			} else {
				run(PojoPath, pathPrefix+"/domain/"+table.JavaBeanName+".java", m)
				run(ControllerPath, pathPrefix+"/controller/"+table.JavaBeanName+"Controller.java", m)
				run(PojoQueryPath, pathPrefix+"/domain/"+table.JavaBeanName+"Query.java", m)
				run(ServicePath, pathPrefix+"/service/"+table.JavaBeanName+"Service.java", m)
				run(ServiceImplPath, pathPrefix+"/service/impl/"+table.JavaBeanName+"ServiceImpl.java", m)
				run(ManagerPath, pathPrefix+"/manager/"+table.JavaBeanName+"Manager.java", m)
				run(ManagerImplPath, pathPrefix+"/manager/impl/"+table.JavaBeanName+"ManagerImpl.java", m)
				run(MapperPath, pathPrefix+"/dao/"+table.JavaBeanName+"Mapper.java", m)
				run(MapperXmlPath, config.ExportPath+"/"+config.ProjectName+"/src/main/resources/mapper/"+table.JavaBeanName+"Mapper.xml", m)
			}
		}

		if config.FrontTheme != "None" {
			err := os.MkdirAll(config.ExportPath+"/"+config.ProjectAdminName+"/src/components/page/"+table.JavaBeanNameLower, os.ModePerm)
			if err != nil {
				panic(err)
			}
			run(vueListPath, config.ExportPath+"/"+config.ProjectAdminName+"/src/components/page/"+table.JavaBeanNameLower+"/list.vue", m)
			//run(vueEditPath, config.ExportPath+"/"+config.ProjectAdminName+"/src/components/page/"+table.JavaBeanNameLower+"/edit.vue", m)
			//run(vueDetailPath, config.ExportPath+"/"+config.ProjectAdminName+"/src/components/page/"+table.JavaBeanNameLower+"/detail.vue", m)
		}

	}
	if config.FrontTheme != "None" {
		m := make(map[string]interface{})
		m["tables"] = schema.Tables
		m["l"] = "{{"
		m["r"] = "}}"
		run(vueIndexPath, config.ExportPath+"/"+config.ProjectAdminName+"/src/router/index.js", m)
		run(vueSidebarPath, config.ExportPath+"/"+config.ProjectAdminName+"/src/components/common/Sidebar.vue", m)
	}

}
