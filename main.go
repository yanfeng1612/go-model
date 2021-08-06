package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yanfeng1612/go-model/utils"
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
	SingleQueryList []string `json:"singleQueryList"`
	InListQuery     []string `json:"inListQuery"`
	UpdateFieldList []string `json:"updateFieldList"`
}

type Theme struct {
	ThemeName string
	templates []string
}

const (
	TemplateRootPath            = "template"                                            // 模板根路径
	TemplateJavaRootPath        = TemplateRootPath + "/java"                            // java模板根路径
	SimpleThemeTemplateRootPath = TemplateJavaRootPath + "/simple/"                     // java模板根路径
	TemplateVueRootPath         = TemplateRootPath + "/vue"                             // java模板根路径
	PomPath                     = TemplateJavaRootPath + "/utils/pom.vm"                // pom路径
	BasicResultPath             = TemplateJavaRootPath + "/utils/BasicResult.vm"        // BasicResultPath路径
	GenericResultPath           = TemplateJavaRootPath + "/utils/GenericResult.vm"      // GenericResult路径
	CodeEnumPath                = TemplateJavaRootPath + "/utils/CodeEnum.vm"           // CodeEnumPath路径
	QueryPath                   = TemplateJavaRootPath + "/utils/Query.vm"              // Query路径
	PageQueryPath               = TemplateJavaRootPath + "/utils/PageQuery.vm"          // PageQuery路径
	IdPageQueryPath             = TemplateJavaRootPath + "/utils/IdPageQuery.vm"        // IdPageQuery路径
	PagenationPath              = TemplateJavaRootPath + "/utils/Pagenation.vm"         // Pagenation路径
	PageQueryWrapperPath        = TemplateJavaRootPath + "/utils/PageQueryWrapper.vm"   // PageQueryWrapper路径
	ListResultPath              = TemplateJavaRootPath + "/utils/ListResult.vm"         // ListResult路径
	PageListResultPath          = TemplateJavaRootPath + "/utils/PageListResult.vm"     // PageListResult路径
	APIEmRequestStatusPath      = TemplateJavaRootPath + "/utils/APIEmRequestStatus.vm" // APIEmRequestStatus路径
	APIMsgCodePath              = TemplateJavaRootPath + "/utils/APIMsgCode.vm"         // APIMsgCode路径
	ResponsePath                = TemplateJavaRootPath + "/utils/Response.vm"           // Response路径
	ResponseTemplatePath        = TemplateJavaRootPath + "/utils/ResponseTemplate.vm"   // ResponseTemplate路径
	CodeConverterPath           = TemplateJavaRootPath + "/utils/CodeConverter.vm"      // CodeConverter路径
	BootstrapPath               = TemplateJavaRootPath + "/utils/Bootstrap.vm"          // Bootstrap路径
	ApplicationConfigPath       = TemplateJavaRootPath + "/utils/application.vm"        // application路径
	GitIgnorePath               = TemplateJavaRootPath + "/utils/gitignore.vm"          // gitignore路径
	AssemblyPath                = TemplateJavaRootPath + "/utils/assembly.vm"           // assembly路径
	StartPath                   = TemplateJavaRootPath + "/utils/start.sh"              // start路径
	StopPath                    = TemplateJavaRootPath + "/utils/stop.sh"               // stop路径

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

var (
	SIMPLE_THEME     = []string{"simple-pojo", "simple-pojoQuery", "simple-controller", "simple-service", "simple-mapper", "simple-mapperXml"}
	STANDARD_THEME   = []string{"pojo", "pojoQuery", "controller", "service", "serviceImpl", "manager", "managerImpl", "mapper", "mapperXml"}
	UtilTemplateList = []string{"basicResult", "genericResult", "codeEnum", "query", "pageQuery", "pageQueryWrapper", "idPageQuery", "pagenation", "listResult", "pageListResult", "aPIEmRequestStatus", "aPIMsgCode", "response", "responseTemplate", "codeConverter", "gitIgnore", "assembly", "start", "stop"}
)

type TemplateObject struct {
	name         string
	templatePath string
	outputPath   string
}

func main() {
	fmt.Println("start generate work")
	args := os.Args
	var configPath string
	if len(args) > 1 {
		configPath = args[1]
	} else {
		configPath = "demo.json"
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
	_, err = json.Marshal(config)
	if err != nil {
		panic(err)
	}

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
				c.UpdateField = 1
			} else {
				c.UpdateField = 0
			}
		}
	}

	templateMap, templateList := prepare(*config, config.BackTheme)

	buildUtils(*config, templateMap)
	buildDynamic(*config, schema, templateList)
	fmt.Println("end generate work")
}

// 初始化
func prepare(config Config, theme string) (map[string]TemplateObject, []TemplateObject) {
	prefix := config.ExportPath + "/" + config.ProjectName + "/src/main/java/" + strings.ReplaceAll(config.PackageName, ".", "/")

	m := make(map[string]TemplateObject)
	// 构建基础模板
	m["pom"] = TemplateObject{"pom", PomPath, config.ExportPath + "/" + config.ProjectName + "/pom.xml"}
	m["bootstrap"] = TemplateObject{"bootstrap", BootstrapPath, prefix + "/Bootstrap.java"}
	m["application"] = TemplateObject{"application", ApplicationConfigPath, config.ExportPath + "/" + config.ProjectName + "/src/main/resources/application.yml"}
	m["gitIgnore"] = TemplateObject{"gitIgnore", GitIgnorePath, config.ExportPath + "/" + config.ProjectName + "/.gitignore"}
	m["assembly"] = TemplateObject{"assembly", AssemblyPath, config.ExportPath + "/" + config.ProjectName + "/src/main/resources/assembly.xml"}
	m["start"] = TemplateObject{"start", StartPath, config.ExportPath + "/" + config.ProjectName + "/src/main/bin/start.sh"}
	m["stop"] = TemplateObject{"stop", StopPath, config.ExportPath + "/" + config.ProjectName + "/src/main/bin/stop.sh"}

	// 工具
	m["basicResult"] = TemplateObject{"basicResult", BasicResultPath, prefix + "/utils/BasicResult.java"}
	m["genericResult"] = TemplateObject{"genericResult", GenericResultPath, prefix + "/utils/GenericResult.java"}
	m["codeEnum"] = TemplateObject{"codeEnum", CodeEnumPath, prefix + "/utils/CodeEnum.java"}
	m["query"] = TemplateObject{"query", QueryPath, prefix + "/utils/Query.java"}
	m["pageQuery"] = TemplateObject{"basicResult", PageQueryPath, prefix + "/utils/PageQuery.java"}
	m["pageQueryWrapper"] = TemplateObject{"pageQueryWrapper", PageQueryWrapperPath, prefix + "/utils/PageQueryWrapper.java"}
	m["idPageQuery"] = TemplateObject{"basicResult", IdPageQueryPath, prefix + "/utils/IdPageQuery.java"}
	m["pagenation"] = TemplateObject{"pagenation", PagenationPath, prefix + "/utils/Pagenation.java"}
	m["listResult"] = TemplateObject{"listResult", ListResultPath, prefix + "/utils/ListResult.java"}
	m["pageListResult"] = TemplateObject{"pageListResult", PageListResultPath, prefix + "/utils/PageListResult.java"}
	m["aPIEmRequestStatus"] = TemplateObject{"aPIEmRequestStatus", APIEmRequestStatusPath, prefix + "/utils/APIEmRequestStatus.java"}
	m["aPIMsgCode"] = TemplateObject{"aPIMsgCode", APIMsgCodePath, prefix + "/utils/APIMsgCode.java"}
	m["response"] = TemplateObject{"response", ResponsePath, prefix + "/utils/Response.java"}
	m["responseTemplate"] = TemplateObject{"responseTemplate", ResponseTemplatePath, prefix + "/utils/ResponseTemplate.java"}
	m["codeConverter"] = TemplateObject{"codeConverter", CodeConverterPath, prefix + "/utils/CodeConverter.java"}

	m["pojo"] = TemplateObject{"pojo", PojoPath, prefix + "/domain/%s.java"}
	m["pojoQuery"] = TemplateObject{"pojoQuery", PojoQueryPath, prefix + "/domain/%sQuery.java"}
	m["controller"] = TemplateObject{"controller", ControllerPath, prefix + "/controller/%sController.java"}
	m["service"] = TemplateObject{"service", ServicePath, prefix + "/service/%sService.java"}
	m["serviceImpl"] = TemplateObject{"serviceImpl", ServiceImplPath, prefix + "/service/impl/%sServiceImpl.java"}
	m["manager"] = TemplateObject{"manager", ManagerPath, prefix + "/manager/%sManager.java"}
	m["managerImpl"] = TemplateObject{"managerImpl", ManagerImplPath, prefix + "/manager/impl/%sManagerImpl.java"}
	m["mapper"] = TemplateObject{"mapper", MapperPath, prefix + "/mapper/%sMapper.java"}
	m["mapperXml"] = TemplateObject{"mapperXml", MapperXmlPath, config.ExportPath + "/" + config.ProjectName + "/src/main/resources/mapper/%sMapper.xml"}

	m["simple-pojo"] = TemplateObject{"pojo", SimpleThemeTemplateRootPath + "pojo.vm", prefix + "/domain/%s.java"}
	m["simple-pojoQuery"] = TemplateObject{"pojoQuery", SimpleThemeTemplateRootPath + "pojoQuery.vm", prefix + "/domain/%sQuery.java"}
	m["simple-controller"] = TemplateObject{"controller", SimpleThemeTemplateRootPath + "controller.vm", prefix + "/controller/%sController.java"}
	m["simple-service"] = TemplateObject{"service", SimpleThemeTemplateRootPath + "service.vm", prefix + "/service/%sService.java"}
	m["simple-mapper"] = TemplateObject{"mapper", SimpleThemeTemplateRootPath + "mapper.vm", prefix + "/mapper/%sMapper.java"}
	m["simple-mapperXml"] = TemplateObject{"mapperXml", SimpleThemeTemplateRootPath + "mapperXml.vm", config.ExportPath + "/" + config.ProjectName + "/src/main/resources/mapper/%sMapper.xml"}

	arr := make([]TemplateObject, 0)

	if theme == "Simple" {
		for _, tem := range SIMPLE_THEME {
			arr = append(arr, m[tem])
		}
	} else if theme == "Standard" {
		for _, tem := range STANDARD_THEME {
			arr = append(arr, m[tem])
		}
	}
	return m, arr
}

func buildUtils(config Config, templateMap map[string]TemplateObject) {
	m := buildCommonMap(config)

	if config.BackTheme == "None" {
		return
	}

	if config.FirstInit {
		pomObj := templateMap["pom"]
		bootstrapObj := templateMap["bootstrap"]
		applicationObj := templateMap["application"]

		run(pomObj.templatePath, pomObj.outputPath, m)
		run(bootstrapObj.templatePath, bootstrapObj.outputPath, m)
		run(applicationObj.templatePath, applicationObj.outputPath, m)
	}

	for _, templateName := range UtilTemplateList {
		template := templateMap[templateName]
		run(template.templatePath, template.outputPath, m)
	}
}

func run(templatePath, outputPath string, m map[string]interface{}) {
	t := buildCommon(templatePath)
	tmpFilePath := outputPath + "Tmp"
	tmpFile, err := createFile(tmpFilePath)
	if err != nil {
		panic(err)
	}
	err = t.Execute(tmpFile, m)
	if err != nil {
		log.Panic("tmpFilePath"+tmpFilePath, err)
	}

	err = tmpFile.Close()
	if err != nil {
		panic(err)
	}

	os.Remove(outputPath)
	tmpFile, _ = os.Open(tmpFilePath)
	file, err := createFile(outputPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReader(tmpFile)
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

	err = tmpFile.Close()
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

func buildDynamic(config Config, schema *utils.Schema, templateObjectList []TemplateObject) {
	for _, table := range schema.Tables {
		m := buildCommonMap(config)
		m["backendObject"] = table
		m["columnSize"] = len(table.Columns)
		m["l"] = "{{"
		m["r"] = "}}"
		m["l1"] = "{"

		for _, tem := range templateObjectList {
			run(tem.templatePath, strings.ReplaceAll(tem.outputPath, "%s", table.JavaBeanName), m)
		}
	}
	//if config.FrontTheme != "None" {
	//	m := make(map[string]interface{})
	//	m["tables"] = schema.Tables
	//	m["l"] = "{{"
	//	m["r"] = "}}"
	//	run(vueIndexPath, config.ExportPath+"/"+config.ProjectAdminName+"/src/router/index.js", m)
	//	run(vueSidebarPath, config.ExportPath+"/"+config.ProjectAdminName+"/src/components/common/Sidebar.vue", m)
	//}
}

func createFile(path string) (*os.File, error) {
	file, _ := os.OpenFile(path, os.O_RDWR, 0777)
	if file != nil {
		return file, nil
	}
	dir, _ := filepath.Split(path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatal("创建目录失败", err)
		return nil, err
	}
	file, _ = os.Create(path)
	return file, nil
}
