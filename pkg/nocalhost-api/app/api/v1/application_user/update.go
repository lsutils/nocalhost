/*
* Copyright (C) 2021 THL A29 Limited, a Tencent company.  All rights reserved.
* This source code is licensed under the Apache License Version 2.0.
 */

package application_user

import (
	"github.com/gin-gonic/gin"
	"github.com/lsutils/nocalhost/internal/nocalhost-api/service"
	"github.com/lsutils/nocalhost/pkg/nocalhost-api/app/api"
	"github.com/lsutils/nocalhost/pkg/nocalhost-api/pkg/errno"
	"github.com/lsutils/nocalhost/pkg/nocalhost-api/pkg/log"
	"github.com/spf13/cast"
)

// BatchInsert batch insert application_user
// only admin user can request this interface
func BatchInsert(c *gin.Context) {
	// userId, _ := c.Get("userId")
	applicationId := cast.ToUint64(c.Param("id"))

	var req ApplicationUsersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warnf("BatchInsert bind params err: %v", err)
		api.SendResponse(c, errno.ErrBind, nil)
		return
	}

	var users []uint64
	for _, user := range req.Users {
		users = append(users, uint64(user))
	}

	err := service.Svc.ApplicationUserSvc.BatchInsert(c, applicationId, users)

	if err != nil {
		api.SendResponse(c, errno.ErrInsertApplicationUser, nil)
		return
	}
	api.SendResponse(c, nil, nil)
}

// BatchDelete batch delete application_user
// only admin user can request this interface
func BatchDelete(c *gin.Context) {
	// userId, _ := c.Get("userId")
	applicationId := cast.ToUint64(c.Param("id"))

	var req ApplicationUsersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warnf("BatchDelete bind params err: %v", err)
		api.SendResponse(c, errno.ErrBind, nil)
		return
	}

	var users []uint64
	for _, user := range req.Users {
		users = append(users, uint64(user))
	}

	err := service.Svc.ApplicationUserSvc.BatchDelete(c, applicationId, users)

	if err != nil {
		api.SendResponse(c, errno.ErrDeleteApplicationUser, nil)
		return
	}
	api.SendResponse(c, nil, nil)
}
