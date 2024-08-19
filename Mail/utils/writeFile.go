package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
	"mail/models" 
)

func WriteNotificationToFile(nofication models.SendNoficitaion) error {
	fileName := fmt.Sprintf("notifications_%s.json", time.Now().Format("20060102_150405"))
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("dosya oluşturulurken hata: %v", err)
	}
	defer file.Close()


	data, err := json.MarshalIndent(nofication, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON oluşturulurken hata: %v", err)
	}

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("dosyaya yazılırken hata: %v", err)
	}

	_, err = file.WriteString("\n")
	if err != nil {
		return fmt.Errorf("dosyaya yeni satır eklenirken hata: %v", err)
	}

	return nil
}
