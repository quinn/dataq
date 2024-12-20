package routes

func (r *Routes) {{ .funcName }}(c echo.Context) error {
    return c.Redirect(http.StatusFound, "/")
}
