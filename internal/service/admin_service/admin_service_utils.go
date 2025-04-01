package admin_service

import "strings"

func (s *adminService) isAdmin(id uint64) bool {
	userID, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return false
	}

	if userID.Role != "admin" {
		return false
	}

	return true
}

func (s *adminService) isValidField(fieldName string) bool {
	allowedFields := []string{"stt", "tts", "ttt", "ls"}

	for _, field := range allowedFields {
		if strings.EqualFold(field, fieldName) {
			return true
		}
	}
	return false
}

func (s *adminService) convertTypeField(fieldType string) string {
	switch fieldType {
	case "stt":
		return "stt_model"
	case "ttt":
		return "ttt_model"
	case "tts":
		return "tts_model"
	case "ls":
		return "ls_model"
	default:
		return ""
	}
}
