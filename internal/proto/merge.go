package proto

import (
	"log"
	"os"
	"regexp"
	"strings"
)

func MergeFile(baseFile string, targetFile string) {
	baseProtoContent, err := os.ReadFile(baseFile)
	if err != nil {
		log.Fatalf("Failed to read %v: %v", baseFile, err)
	}

	// 读取 b.proto 文件
	targetProtoContent, err := os.ReadFile(targetFile)
	if err != nil {
		log.Fatalf("Failed to read %v: %v", targetFile, err)
	}

	// 更新正则表达式以确保更稳健的匹配
	messagePattern := `message\s+\w+\s+{[^}]*}`
	servicePattern := `service\s+\w+\s+{[^}]*}`

	// 使用正则表达式提取 message 和 service
	messageRegex := regexp.MustCompile(messagePattern)
	serviceRegex := regexp.MustCompile(servicePattern)

	messages := messageRegex.FindAllString(string(baseProtoContent), -1)
	services := serviceRegex.FindAllString(string(baseProtoContent), -1)

	// 如果没有找到 message 或 service，输出提示
	if len(messages) == 0 && len(services) == 0 {
		log.Printf("No message or service found in %s.", targetFile)
	}

	// 提取 base.proto 中现有的 message 和 service 部分
	bProtoStr := string(targetProtoContent)
	bMessagesPattern := `message\s+\w+\s+{[^}]*}`
	bServicesPattern := `service\s+\w+\s+{[^}]*}`

	bMessagesRegex := regexp.MustCompile(bMessagesPattern)
	bServicesRegex := regexp.MustCompile(bServicesPattern)

	// 提取 targetFile 中现有的 messages 和 services
	bMessages := bMessagesRegex.FindAllString(bProtoStr, -1)
	bServices := bServicesRegex.FindAllString(bProtoStr, -1)

	// 去重：检查 base.proto 中的 message 和 service 是否已经存在于 target.proto 中
	var contentToAppend string
	for _, message := range messages {
		// 只追加 target.proto 中没有的 message
		if !contains(bMessages, message) {
			contentToAppend += "\n\n" + message
		}
	}

	var servicesToAppend string
	for _, service := range services {
		if !contains(bServices, service) {
			servicesToAppend += "\n\n" + service
		}
	}

	// 确保 message 出现在 target.proto 的 message 部分之后
	bProtoUpdated := ""
	if len(bMessages) > 0 {
		// 找到最后一个 message 部分，将新的 message 追加到这个位置
		lastMessageEnd := strings.LastIndex(bProtoStr, bMessages[len(bMessages)-1])
		bProtoUpdated = bProtoStr[:lastMessageEnd+len(bMessages[len(bMessages)-1])] + contentToAppend + bProtoStr[lastMessageEnd+len(bMessages[len(bMessages)-1]):]
	} else {
		// 如果target.proto 中没有 message，直接将 message 内容添加到末尾
		bProtoUpdated = bProtoStr + "\n\n" + contentToAppend
	}

	// 追加 service 内容
	bProtoUpdated += servicesToAppend

	// 写回到 target.proto 文件
	err = os.WriteFile(targetFile, []byte(bProtoUpdated), 0644)
	if err != nil {
		log.Fatalf("Failed to write b.proto: %v", err)
	}

}
func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
