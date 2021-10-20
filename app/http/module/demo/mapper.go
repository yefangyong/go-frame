package demo

import "github.com/yefangyong/go-frame/app/provider/demo"

func UserModelsToUserDTOs(models []UserModel) []UserDTO {
	var ret []UserDTO
	for _, model := range models {
		t := UserDTO{Id: model.UserId, Name: model.Name}
		ret = append(ret, t)
	}
	return ret
}

func StudentsToUserDTOs(students []demo.Student) []UserDTO {
	var ret []UserDTO
	for _, student := range students {
		t := UserDTO{Id: student.ID, Name: student.Name}
		ret = append(ret, t)
	}
	return ret
}
