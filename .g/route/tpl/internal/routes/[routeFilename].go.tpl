package routes

func {{ .funcName }}(c echo.Context) error {
    return c.Redirect(http.StatusFound, "/")
}
