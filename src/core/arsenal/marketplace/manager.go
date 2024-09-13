package marketplace

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/olekukonko/tablewriter"
)

// // ManageTool exécute une ou plusieurs commandes pour une action donnée (install, remove, update) d'un outil en fonction du gestionnaire de paquets.
// func ManageTool(tool Tool, packageManager, actionType string) error {
// 	if pmActions, exists := tool.Actions[packageManager]; exists {
// 		var commands []string
// 		switch actionType {
// 		case "install":
// 			commands = pmActions.Install
// 		case "remove":
// 			commands = pmActions.Remove
// 		case "update":
// 			commands = pmActions.Update
// 		default:
// 			return fmt.Errorf("unknown action type: %s", actionType)
// 		}

// 		if len(commands) == 0 {
// 			return fmt.Errorf("no %s actions defined for tool: %s with package manager: %s", actionType, tool.Name, packageManager)
// 		}

// 		for _, cmd := range commands {
// 			if err := runCommand(cmd, tool.Name, actionType); err != nil {
// 				return err
// 			}
// 		}
// 		return nil
// 	}
// 	return fmt.Errorf("no actions found for package manager: %s", packageManager)
// }

// InstallTool exécute toutes les commandes d'installation d'un outil en fonction du gestionnaire de paquets.
func InstallTool(tool Tool, packageManager string) error {
	if pmActions, exists := tool.Actions[packageManager]; exists {
		for _, cmd := range pmActions.Install {
			if err := runCommand(cmd, tool.Name, "install"); err != nil {
				return err
			}
		}
		return nil
	}
	return fmt.Errorf("no install actions defined for tool: %s with package manager: %s", tool.Name, packageManager)
}

// RemoveTool exécute toutes les commandes de suppression d'un outil en fonction du gestionnaire de paquets.
func RemoveTool(tool Tool, packageManager string) error {
	if pmActions, exists := tool.Actions[packageManager]; exists {
		for _, cmd := range pmActions.Remove {
			if err := runCommand(cmd, tool.Name, "remove"); err != nil {
				return err
			}
		}
		return nil
	}
	return fmt.Errorf("no remove actions defined for tool: %s with package manager: %s", tool.Name, packageManager)
}

// UpdateTool exécute toutes les commandes de mise à jour d'un outil en fonction du gestionnaire de paquets.
func UpdateTool(tool Tool, packageManager string) error {
	if pmActions, exists := tool.Actions[packageManager]; exists {
		for _, cmd := range pmActions.Update {
			if err := runCommand(cmd, tool.Name, "update"); err != nil {
				return err
			}
		}
		return nil
	}
	return fmt.Errorf("no update actions defined for tool: %s with package manager: %s", tool.Name, packageManager)
}

// TestTool teste si l'outil est installé en exécutant les commandes de test définies.
func TestTool(tool Tool) error {
	for _, command := range tool.TestCommands {
		cmd := exec.Command("sh", "-c", command) // Exécute la commande dans un shell
		output, err := cmd.CombinedOutput()      // Combine stdout et stderr
		if err != nil {
			return fmt.Errorf("error running test command for tool %s: %s\nOutput: %s", tool.Name, err, output)
		}

		// Vérifie si la commande a produit une sortie et l'affiche pour le débogage
		fmt.Printf("Test command for tool %s:\n%s\n", tool.Name, output)
	}
	return nil
}

// runCommand exécute une commande shell spécifique pour un outil et une action donnée.
func runCommand(cmdStr, toolName, action string) error {
	fmt.Printf("Running %s command for tool %s: %s\n", action, toolName, cmdStr)
	// Exécution de la commande shell
	cmd := exec.Command("sh", "-c", cmdStr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running %s command for tool %s: %v\nOutput: %s", action, toolName, err, string(output))
	}
	fmt.Printf("Command output: %s\n", string(output))
	return nil
}

// ShowTools affiche les informations sur les outils disponibles dans un format tabulaire.
func ShowTools(tools []Tool) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Description", "Version", "Author", "Commands", "Categories", "Tags", "Metadata"})

	if len(tools) == 0 {
		fmt.Println("Aucun outil trouvé.")
		return
	}

	for _, tool := range tools {
		commands := ""
		for _, cmd := range tool.TestCommands {
			commands += cmd + "; "
		}

		metadata := fmt.Sprintf("License: %s, Source: %s, Logo: %s", tool.Metadata.License, tool.Metadata.Source, tool.Metadata.Logo_url)

		// Ajouter les informations de l'outil dans la table
		table.Append([]string{
			tool.ID,
			tool.Name,
			tool.Description,
			tool.Version,
			tool.Author,
			commands,
			fmt.Sprintf("%v", tool.Categories),
			fmt.Sprintf("%v", tool.Tags),
			metadata,
		})
	}

	// Afficher la table
	table.Render()
}

// FindToolsByCategory recherche les outils par catégorie dans le registry.
func (r *Registry) FindToolsByCategory(category string) []Tool {
	var results []Tool
	for _, tool := range r.Tools {
		for _, cat := range tool.Categories {
			if cat == category {
				results = append(results, tool)
				break
			}
		}
	}
	return results
}

// FindToolsByTag recherche les outils par tag
func (r *Registry) FindToolsByTag(tag string) []Tool {
	var results []Tool
	for _, tool := range r.Tools {
		for _, t := range tool.Tags {
			if t == tag {
				results = append(results, tool)
				break
			}
		}
	}
	return results
}
