package storage

import (
	"fmt"
	"strings"
	"time"

	"github.com/sebasvalencia/clical/pkg/calendar"
	"github.com/sebasvalencia/clical/pkg/user"
)

// entryToMarkdown convierte una Entry a formato Markdown
func entryToMarkdown(entry *calendar.Entry) string {
	var md strings.Builder

	// Título
	md.WriteString(fmt.Sprintf("# %s\n\n", entry.Title))

	// Metadata principal
	md.WriteString(fmt.Sprintf("**Fecha:** %s  \n", entry.DateTime.Format("2006-01-02")))
	md.WriteString(fmt.Sprintf("**Hora:** %s  \n", entry.DateTime.Format("15:04")))
	md.WriteString(fmt.Sprintf("**Duración:** %d minutos  \n", entry.Duration))

	if entry.Location != "" {
		md.WriteString(fmt.Sprintf("**Ubicación:** %s  \n", entry.Location))
	}

	if len(entry.Tags) > 0 {
		tags := make([]string, len(entry.Tags))
		for i, tag := range entry.Tags {
			tags[i] = "#" + tag
		}
		md.WriteString(fmt.Sprintf("**Tags:** %s  \n", strings.Join(tags, " ")))
	}

	md.WriteString("\n")

	// Notas
	if entry.Notes != "" {
		md.WriteString("## Notas\n\n")
		md.WriteString(entry.Notes)
		md.WriteString("\n\n")
	}

	// Metadata adicional
	if len(entry.Metadata) > 0 {
		md.WriteString("## Metadata\n\n")
		for k, v := range entry.Metadata {
			md.WriteString(fmt.Sprintf("- **%s:** %s\n", k, v))
		}
		md.WriteString("\n")
	}

	// Footer con timestamps
	md.WriteString("---\n\n")
	md.WriteString(fmt.Sprintf("*Creado: %s*  \n", entry.CreatedAt.Format("2006-01-02 15:04")))
	md.WriteString(fmt.Sprintf("*Actualizado: %s*  \n", entry.UpdatedAt.Format("2006-01-02 15:04")))
	md.WriteString(fmt.Sprintf("*ID: %s*\n", entry.ID))

	return md.String()
}

// userToMarkdown convierte un User a formato Markdown
func userToMarkdown(u *user.User) string {
	var md strings.Builder

	// Título
	md.WriteString(fmt.Sprintf("# Usuario: %s\n\n", u.Name))

	// Información básica
	md.WriteString(fmt.Sprintf("**ID:** %s  \n", u.ID))
	md.WriteString(fmt.Sprintf("**Timezone:** %s  \n", u.Timezone))
	md.WriteString(fmt.Sprintf("**Miembro desde:** %s\n\n", u.Created.Format("2006-01-02")))

	// Configuración
	md.WriteString("## Configuración\n\n")
	md.WriteString(fmt.Sprintf("- **Duración por defecto:** %d minutos\n", u.Config.DefaultDuration))
	md.WriteString(fmt.Sprintf("- **Formato de fecha:** %s\n", u.Config.DateFormat))
	md.WriteString(fmt.Sprintf("- **Formato de hora:** %s\n", u.Config.TimeFormat))

	firstDay := "Domingo"
	if u.Config.FirstDayOfWeek == 1 {
		firstDay = "Lunes"
	}
	md.WriteString(fmt.Sprintf("- **Primer día de semana:** %s\n", firstDay))

	md.WriteString("\n---\n\n")
	md.WriteString(fmt.Sprintf("*Última actualización: %s*\n", time.Now().Format("2006-01-02")))

	return md.String()
}
