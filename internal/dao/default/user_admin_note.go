// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"unibee/internal/dao/default/internal"
)

// internalUserAdminNoteDao is internal type for wrapping internal DAO implements.
type internalUserAdminNoteDao = *internal.UserAdminNoteDao

// userAdminNoteDao is the data access object for table user_admin_note.
// You can define custom methods on it to extend its functionality as you wish.
type userAdminNoteDao struct {
	internalUserAdminNoteDao
}

var (
	// UserAdminNote is globally public accessible object for table user_admin_note operations.
	UserAdminNote = userAdminNoteDao{
		internal.NewUserAdminNoteDao(),
	}
)

// Fill with you ideas below.