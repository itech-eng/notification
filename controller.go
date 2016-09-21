package notification

import (
	"net/http"

	"github.com/qor/admin"
	"github.com/qor/responder"
)

type controller struct {
	Notification *Notification
	action       *Action
}

func (c *controller) List(context *admin.Context) {
	context.Set("Notification", c.Notification)
	context.Execute("notifications/notifications", map[string]interface{}{
		"Messages": c.Notification.GetNotifications(context.CurrentUser, context.Context),
	})
}

func (c *controller) Action(context *admin.Context) {
	action := c.action
	message := c.Notification.GetNotification(context.CurrentUser, context.ResourceID, context.Context)
	context.Set("Notification", c.Notification)

	if context.Request.Method == "GET" {
		context.Execute("action", action)
	} else {
		var actionArgument = &ActionArgument{
			Message: message,
			Context: context,
		}

		if action.Resource != nil {
			result := action.Resource.NewStruct()
			action.Resource.Decode(context.Context, result)
			actionArgument.Argument = result
		}

		if err := action.Handle(actionArgument); err == nil {
			flash := string(context.Admin.T(context.Context, "qor_admin.actions.executed_successfully", "Action {{.Name}}: Executed successfully", action))
			responder.With("html", func() {
				context.Flash(flash, "success")
				http.Redirect(context.Writer, context.Request, context.Request.Referer(), http.StatusFound)
			}).With("json", func() {
				notification := c.Notification.GetNotification(context.CurrentUser, context.ResourceID, context.Context)
				context.JSON("OK", map[string]string{"status": "ok", "message": flash, "notification": string(context.Render("notification", notification))})
			}).Respond(context.Request)
		} else {
			notification := c.Notification.GetNotification(context.CurrentUser, context.ResourceID, context.Context)
			flash := string(context.Admin.T(context.Context, "qor_admin.actions.executed_failed", "Action {{.Name}}: Failed to execute", action))
			context.JSON("OK", map[string]string{"status": "error", "error": flash, "notification": string(context.Render("notification", notification))})
		}
	}
}

func (c *controller) UndoAction(context *admin.Context) {
	action := c.action
	message := c.Notification.GetNotification(context.CurrentUser, context.ResourceID, context.Context)
	context.Set("Notification", c.Notification)

	var actionArgument = &ActionArgument{
		Message: message,
		Context: context,
	}

	if action.Resource != nil {
		result := action.Resource.NewStruct()
		action.Resource.Decode(context.Context, result)
		actionArgument.Argument = result
	}

	if err := action.Undo(actionArgument); err == nil {
		flash := string(context.Admin.T(context.Context, "qor_admin.actions.executed_successfully", "Action {{.Name}}: Undoed successfully", action))
		responder.With("html", func() {
			context.Flash(flash, "success")
			http.Redirect(context.Writer, context.Request, context.Request.Referer(), http.StatusFound)
		}).With("json", func() {
			notification := c.Notification.GetNotification(context.CurrentUser, context.ResourceID, context.Context)
			context.JSON("OK", map[string]string{"status": "ok", "message": flash, "notification": string(context.Render("notification", notification))})
		}).Respond(context.Request)
	} else {
		notification := c.Notification.GetNotification(context.CurrentUser, context.ResourceID, context.Context)
		flash := string(context.Admin.T(context.Context, "qor_admin.actions.executed_failed", "Action {{.Name}}: Failed to undo", action))
		context.JSON("OK", map[string]string{"status": "error", "error": flash, "notification": string(context.Render("notification", notification))})
	}
}
