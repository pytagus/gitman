package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
)

type GitManager struct {
	currentPath string
	scanner     *bufio.Scanner
}

func NewGitManager() *GitManager {
	currentPath, _ := os.Getwd()
	return &GitManager{
		currentPath: currentPath,
		scanner:     bufio.NewScanner(os.Stdin),
	}
}

// General Helpers
func (gm *GitManager) getUserInput() string {
	if gm.scanner.Scan() {
		return strings.TrimSpace(gm.scanner.Text())
	}
	if err := gm.scanner.Err(); err != nil {
		fmt.Printf("%sErreur de lecture: %v%s\n", ColorRed, err, ColorReset)
	}
	return ""
}

func (gm *GitManager) pause() {
	fmt.Printf("\n%sAppuyez sur Entrée pour continuer...%s", ColorYellow, ColorReset)
	gm.getUserInput()
}

func (gm *GitManager) clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		// Fallback for windows
		cmd = exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func (gm *GitManager) runGitCommand(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = gm.currentPath
	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

func (gm *GitManager) isGitRepo() bool {
	_, err := os.Stat(filepath.Join(gm.currentPath, ".git"))
	return err == nil
}

// UI and Menu
func (gm *GitManager) printHeader() {
	fmt.Printf("%s%s╔════════════════════════════════════════════════════════════════╗%s\n", ColorBold, ColorCyan, ColorReset)
	fmt.Printf("%s%s║                     🔧 GIT MANAGER CLI 🔧                     ║%s\n", ColorBold, ColorCyan, ColorReset)
	fmt.Printf("%s%s╚════════════════════════════════════════════════════════════════╝%s\n", ColorBold, ColorCyan, ColorReset)
	fmt.Printf("%sRépertoire actuel: %s%s%s\n\n", ColorYellow, ColorWhite, gm.currentPath, ColorReset)
}

// Modification de la fonction showMenu() pour inclure un accès rapide
func (gm *GitManager) showMenu() {
	gm.printHeader()

	if gm.isGitRepo() {
		gm.showQuickStatus()
		gm.showQuickActions() // Nouvelle fonction pour les accès rapides
	}

	fmt.Printf("%s%s═══════════════════════════════════════════════════════════════════%s\n", ColorBold, ColorBlue, ColorReset)
	fmt.Printf("%s%s                           📋 MENU PRINCIPAL                           %s\n", ColorBold, ColorBlue, ColorReset)
	fmt.Printf("%s%s═══════════════════════════════════════════════════════════════════%s\n", ColorBold, ColorBlue, ColorReset)

	// SECTION ACCÈS RAPIDE
	fmt.Printf("%s%s⚡ ACCÈS RAPIDE:%s\n", ColorBold, ColorYellow, ColorReset)
	fmt.Printf("%s S%s  📊 Statut détaillé\n", ColorCyan, ColorReset)
	fmt.Printf("%s C%s  📦 Commits (nouveau/historique)\n", ColorCyan, ColorReset)
	fmt.Printf("%s F%s  📁 Fichiers (add/diff)\n", ColorCyan, ColorReset)
	fmt.Printf("%s B%s  🌿 Branches (créer/changer)\n", ColorCyan, ColorReset)
	fmt.Printf("%s R%s  🔄 Remote (push/pull)\n", ColorCyan, ColorReset)

	fmt.Printf("\n%s%s📋 MENU COMPLET:%s\n", ColorBold, ColorGreen, ColorReset)
	fmt.Printf("%s 1.%s  📊 Statut détaillé du dépôt\n", ColorGreen, ColorReset)
	fmt.Printf("%s 2.%s  🌿 Gestion des branches\n", ColorGreen, ColorReset)
	fmt.Printf("%s 3.%s  📦 Gestion des commits\n", ColorGreen, ColorReset)
	fmt.Printf("%s 4.%s  🔄 Gestion des remotes\n", ColorGreen, ColorReset)
	fmt.Printf("%s 5.%s  📁 Gestion des fichiers\n", ColorGreen, ColorReset)
	fmt.Printf("%s 6.%s  🏷️  Gestion des tags\n", ColorGreen, ColorReset)
	fmt.Printf("%s 7.%s  🗂️  Gestion des stash\n", ColorGreen, ColorReset)
	fmt.Printf("%s 8.%s  📈 Statistiques et logs\n", ColorGreen, ColorReset)
	fmt.Printf("%s 9.%s  🔧 Outils et configuration\n", ColorGreen, ColorReset)
	fmt.Printf("%s10.%s  📂 Changer de répertoire\n", ColorGreen, ColorReset)
	fmt.Printf("%s11.%s  🚀 Initialiser un nouveau dépôt\n", ColorGreen, ColorReset)
	fmt.Printf("%s 0.%s  🚪 Quitter\n", ColorRed, ColorReset)

	fmt.Printf("\n%sChoisissez une option: %s", ColorYellow, ColorReset)
}

func (gm *GitManager) showQuickStatus() {
	branch := gm.getCurrentBranch()
	status := gm.getGitStatus()

	fmt.Printf("%s%s📍 Branche actuelle: %s%s%s\n", ColorBold, ColorGreen, ColorCyan, branch, ColorReset)

	if status != "" {
		lines := strings.Split(strings.TrimSpace(status), "\n")
		modified := 0
		staged := 0
		untracked := 0

		for _, line := range lines {
			if strings.HasPrefix(line, " M") || strings.HasPrefix(line, "AM") {
				modified++
			} else if strings.HasPrefix(line, "M ") || strings.HasPrefix(line, "A ") {
				staged++
			} else if strings.HasPrefix(line, "??") {
				untracked++
			}
		}

		if staged > 0 {
			fmt.Printf("%s✓ %d fichier(s) en stage%s ", ColorGreen, staged, ColorReset)
		}
		if modified > 0 {
			fmt.Printf("%s⚠ %d fichier(s) modifié(s)%s ", ColorYellow, modified, ColorReset)
		}
		if untracked > 0 {
			fmt.Printf("%s? %d fichier(s) non suivi(s)%s ", ColorRed, untracked, ColorReset)
		}
		fmt.Println()
	}
	fmt.Println()
}

// Nouvelle fonction pour afficher les actions rapides avec le statut
func (gm *GitManager) showQuickActions() {
	fmt.Printf("%s%s⚡ ACTIONS RAPIDES DISPONIBLES:%s\n", ColorBold, ColorYellow, ColorReset)

	// Analyser le contexte pour suggérer des actions
	status := gm.getGitStatus()
	staged, _ := gm.runGitCommand("diff", "--cached", "--name-only")
	currentBranch := gm.getCurrentBranch()

	if status != "" {
		if staged != "" {
			fmt.Printf("%s   💡 Vous avez des fichiers en stage → tapez 'C' pour commiter%s\n", ColorGreen, ColorReset)
		} else {
			fmt.Printf("%s   💡 Fichiers modifiés détectés → tapez 'F' pour les ajouter%s\n", ColorYellow, ColorReset)
		}
	}

	if currentBranch == "main" || currentBranch == "master" {
		fmt.Printf("%s   💡 Sur branche principale → tapez 'B' pour créer une feature branch%s\n", ColorCyan, ColorReset)
	}

	fmt.Println()
}

// Git data retrievers
func (gm *GitManager) getCurrentBranch() string {
	output, err := gm.runGitCommand("branch", "--show-current")
	if err != nil {
		return "unknown"
	}
	return output
}

func (gm *GitManager) getGitStatus() string {
	output, _ := gm.runGitCommand("status", "--porcelain")
	return output
}

// Menu Handlers
// Version améliorée de handleDetailedStatus avec intelligence contextuelle
func (gm *GitManager) handleDetailedStatus() {
	if !gm.isGitRepo() {
		fmt.Printf("%s❌ Ce répertoire n'est pas un dépôt Git!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	fmt.Printf("%s%s📊 STATUT INTELLIGENT DU DÉPÔT%s\n", ColorBold, ColorBlue, ColorReset)
	fmt.Println(strings.Repeat("═", 60))

	// 1. INFORMATIONS DE BASE
	gm.showBasicRepoInfo()

	// 2. ÉTAT DES FICHIERS (avec analyse intelligente)
	gm.showIntelligentFileStatus()

	// 3. INFORMATIONS SUR LES BRANCHES
	gm.showBranchInfo()

	// 4. SYNCHRONISATION REMOTE
	gm.showRemoteSync()

	// 5. DERNIÈRE ACTIVITÉ
	gm.showRecentActivity()

	// 6. SUGGESTIONS INTELLIGENTES
	gm.showIntelligentSuggestions()

	gm.pause()
}

// Informations de base du dépôt
func (gm *GitManager) showBasicRepoInfo() {
	currentBranch := gm.getCurrentBranch()
	lastCommit, _ := gm.runGitCommand("log", "-1", "--pretty=format:%h - %s (%an, %ar)")

	fmt.Printf("%s🏠 DÉPÔT:%s %s\n", ColorBold, ColorReset, filepath.Base(gm.currentPath))
	fmt.Printf("%s🌿 BRANCHE ACTUELLE:%s %s%s%s\n", ColorBold, ColorReset, ColorCyan, currentBranch, ColorReset)
	fmt.Printf("%s📦 DERNIER COMMIT:%s %s\n", ColorBold, ColorReset, lastCommit)

	// Compter les commits
	totalCommits, _ := gm.runGitCommand("rev-list", "--count", "HEAD")
	if totalCommits != "" {
		fmt.Printf("%s📊 TOTAL COMMITS:%s %s\n", ColorBold, ColorReset, totalCommits)
	}

	fmt.Println()
}

// Analyse intelligente de l'état des fichiers
func (gm *GitManager) showIntelligentFileStatus() {
 status := gm.getGitStatus()
 staged, _ := gm.runGitCommand("diff", "--cached", "--name-only")

 fmt.Printf("%s%s📁 ÉTAT DES FICHIERS%s\n", ColorBold, ColorGreen, ColorReset)

 if status == "" && staged == "" {
  fmt.Printf("%s✨ Working directory clean - Aucun changement détecté%s\n", ColorGreen, ColorReset)
  fmt.Println()
  return
 }

 // Analyser et catégoriser les changements
 var modified, added, deleted, renamed, untracked, stagedFiles []string

 if status != "" {
  lines := strings.Split(strings.TrimSpace(status), "\n")
  for _, line := range lines {
   if len(line) < 3 {
    continue
   }
   statusCode := line[:2]
   fileName := line[3:]

   switch statusCode {
   case "M ", "AM":
    stagedFiles = append(stagedFiles, fileName+" (modifié)")
   case " M", "MM":
    modified = append(modified, fileName)
   case "A ":
    stagedFiles = append(stagedFiles, fileName+" (nouveau)")
    added = append(added, fileName) // Utiliser la variable 'added'
   case "D ":
    stagedFiles = append(stagedFiles, fileName+" (supprimé)")
   case " D":
    deleted = append(deleted, fileName)
   case "R ":
    renamed = append(renamed, fileName)
   case "??":
    untracked = append(untracked, fileName)
   }
  }
 }

 // Affichage organisé avec couleurs et statistiques
 if len(stagedFiles) > 0 {
  fmt.Printf("%s✅ FICHIERS EN STAGE (%d):%s\n", ColorGreen, len(stagedFiles), ColorReset)
  for _, file := range stagedFiles {
   fmt.Printf("   %s▶%s %s\n", ColorGreen, ColorReset, file)
  }

  // Statistiques des changements stagés
  statsOutput, _ := gm.runGitCommand("diff", "--cached", "--stat")
  if statsOutput != "" {
   fmt.Printf("%s   📊 Statistiques:%s\n", ColorCyan, ColorReset)
   lines := strings.Split(statsOutput, "\n")
   for _, line := range lines {
    if line != "" && !strings.Contains(line, "file") {
     fmt.Printf("   %s\n", line)
    }
   }
  }
  fmt.Println()
 }

 if len(modified) > 0 {
  fmt.Printf("%s⚠️  FICHIERS MODIFIÉS (%d):%s\n", ColorYellow, len(modified), ColorReset)
  for i, file := range modified {
   if i < 10 { // Limiter l'affichage
    fmt.Printf("   %s●%s %s\n", ColorYellow, ColorReset, file)
   } else if i == 10 {
    fmt.Printf("   %s... et %d autre(s)%s\n", ColorYellow, len(modified)-10, ColorReset)
    break
   }
  }
  fmt.Println()
 }

 // Afficher les fichiers ajoutés s'il y en a (utilisation de la variable 'added')
 if len(added) > 0 {
  fmt.Printf("%s➕ NOUVEAUX FICHIERS (%d):%s\n", ColorBlue, len(added), ColorReset)
  for _, file := range added {
   fmt.Printf("   %s+%s %s\n", ColorBlue, ColorReset, file)
  }
  fmt.Println()
 }

 if len(untracked) > 0 {
  fmt.Printf("%s❓ FICHIERS NON SUIVIS (%d):%s\n", ColorRed, len(untracked), ColorReset)
  for i, file := range untracked {
   if i < 5 { // Limiter l'affichage pour les non suivis
    fmt.Printf("   %s?%s %s\n", ColorRed, ColorReset, file)
   } else if i == 5 {
    fmt.Printf("   %s... et %d autre(s)%s\n", ColorRed, len(untracked)-5, ColorReset)
    break
   }
  }
  fmt.Println()
 }

 if len(deleted) > 0 {
  fmt.Printf("%s🗑️  FICHIERS SUPPRIMÉS (%d):%s\n", ColorRed, len(deleted), ColorReset)
  for _, file := range deleted {
   fmt.Printf("   %s✗%s %s\n", ColorRed, ColorReset, file)
  }
  fmt.Println()
 }

 if len(renamed) > 0 {
  fmt.Printf("%s🔄 FICHIERS RENOMMÉS (%d):%s\n", ColorCyan, len(renamed), ColorReset)
  for _, file := range renamed {
   fmt.Printf("   %s↻%s %s\n", ColorCyan, ColorReset, file)
  }
  fmt.Println()
 }
}

// Informations détaillées sur les branches
func (gm *GitManager) showBranchInfo() {
	fmt.Printf("%s%s🌿 INFORMATIONS BRANCHES%s\n", ColorBold, ColorGreen, ColorReset)

	currentBranch := gm.getCurrentBranch()

	// Branches locales
	localBranches, _ := gm.runGitCommand("branch", "--format=%(refname:short)")
	localCount := 0
	if localBranches != "" {
		localCount = len(strings.Split(strings.TrimSpace(localBranches), "\n"))
	}

	// Branches remote
	remoteBranches, _ := gm.runGitCommand("branch", "-r", "--format=%(refname:short)")
	remoteCount := 0
	if remoteBranches != "" {
		remoteCount = len(strings.Split(strings.TrimSpace(remoteBranches), "\n"))
	}

	fmt.Printf("%s📍 Branche actuelle:%s %s%s%s\n", ColorBlue, ColorReset, ColorCyan, currentBranch, ColorReset)
	fmt.Printf("%s🏠 Branches locales:%s %d\n", ColorBlue, ColorReset, localCount)
	fmt.Printf("%s🌐 Branches remote:%s %d\n", ColorBlue, ColorReset, remoteCount)

	// Vérifier si on est sur main/master
	if currentBranch == "main" || currentBranch == "master" {
		fmt.Printf("%s⚠️  Vous êtes sur la branche principale%s\n", ColorYellow, ColorReset)
	}

	// Dernières branches utilisées
	recentBranches, err := gm.runGitCommand("for-each-ref", "--count=3", "--sort=-committerdate",
		"--format=%(refname:short) (%(committerdate:relative))", "refs/heads/")
	if err == nil && recentBranches != "" {
		fmt.Printf("%s🕐 Branches récemment utilisées:%s\n", ColorCyan, ColorReset)
		lines := strings.Split(recentBranches, "\n")
		for _, line := range lines {
			if line != "" {
				fmt.Printf("   %s\n", line)
			}
		}
	}
	fmt.Println()
}

// État de synchronisation avec les remotes
func (gm *GitManager) showRemoteSync() {
	fmt.Printf("%s%s🔄 SYNCHRONISATION REMOTE%s\n", ColorBold, ColorPurple, ColorReset)

	// Vérifier les remotes configurés
	remotes, _ := gm.runGitCommand("remote", "-v")
	if remotes == "" {
		fmt.Printf("%s❌ Aucun remote configuré%s\n", ColorRed, ColorReset)
		fmt.Println()
		return
	}

	// Extraire le remote principal (généralement origin)
	originURL := ""
	remoteLines := strings.Split(remotes, "\n")
	for _, line := range remoteLines {
		if strings.Contains(line, "origin") && strings.Contains(line, "(fetch)") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				originURL = parts[1]
			}
			break
		}
	}

	if originURL != "" {
		fmt.Printf("%s🌐 Remote principal:%s %s\n", ColorBlue, ColorReset, originURL)
	}

	// Vérifier l'état de synchronisation (commits en avance/retard)
	// Effectuer un fetch silencieux pour s'assurer que les informations sont à jour
	gm.runGitCommand("fetch", "origin")
	ahead, _ := gm.runGitCommand("rev-list", "--count", "@{u}..HEAD")
	behind, _ := gm.runGitCommand("rev-list", "--count", "HEAD..@{u}")

	if ahead != "" && ahead != "0" {
		fmt.Printf("%s📤 Commits à pusher:%s %s commit(s)\n", ColorYellow, ColorReset, ahead)
	}
	if behind != "" && behind != "0" {
		fmt.Printf("%s📥 Commits à récupérer:%s %s commit(s)\n", ColorYellow, ColorReset, behind)
	}

	if (ahead == "" || ahead == "0") && (behind == "" || behind == "0") {
		fmt.Printf("%s✅ Branche synchronisée avec le remote%s\n", ColorGreen, ColorReset)
	}

	// Dernière synchronisation
	lastFetch, _ := gm.runGitCommand("log", "-1", "--pretty=format:%ar", "FETCH_HEAD")
	if lastFetch != "" {
		fmt.Printf("%s🕐 Dernier fetch:%s %s\n", ColorCyan, ColorReset, lastFetch)
	}

	fmt.Println()
}

// Activité récente du dépôt
func (gm *GitManager) showRecentActivity() {
	fmt.Printf("%s%s📈 ACTIVITÉ RÉCENTE%s\n", ColorBold, ColorBlue, ColorReset)

	// Derniers commits (3 derniers)
	recentCommits, _ := gm.runGitCommand("log", "-3", "--pretty=format:%C(yellow)%h%C(reset) %s %C(cyan)(%an, %ar)%C(reset)")
	if recentCommits != "" {
		fmt.Printf("%s📦 Derniers commits:%s\n", ColorGreen, ColorReset)
		lines := strings.Split(recentCommits, "\n")
		for _, line := range lines {
			if line != "" {
				fmt.Printf("   %s\n", line)
			}
		}
	}

	// Activité des contributeurs (si plusieurs contributeurs)
	contributors, _ := gm.runGitCommand("shortlog", "-sn", "--since=1.week.ago")
	if contributors != "" {
		contributorLines := strings.Split(strings.TrimSpace(contributors), "\n")
		if len(contributorLines) > 1 {
			fmt.Printf("%s👥 Activité cette semaine:%s\n", ColorCyan, ColorReset)
			for i, line := range contributorLines {
				if i < 3 && line != "" { // Top 3 contributeurs
					fmt.Printf("   %s\n", line)
				}
			}
		}
	}

	fmt.Println()
}

// Suggestions intelligentes basées sur l'état du dépôt
func (gm *GitManager) showIntelligentSuggestions() {
	fmt.Printf("%s%s💡 SUGGESTIONS INTELLIGENTES%s\n", ColorBold, ColorYellow, ColorReset)

	status := gm.getGitStatus()
	staged, _ := gm.runGitCommand("diff", "--cached", "--name-only")
	currentBranch := gm.getCurrentBranch()
	ahead, _ := gm.runGitCommand("rev-list", "--count", "@{u}..HEAD")
	behind, _ := gm.runGitCommand("rev-list", "--count", "HEAD..@{u}")

	suggestions := []string{}

	// Suggestions basées sur l'état des fichiers
	if status != "" && staged == "" {
		suggestions = append(suggestions, "📁 Des fichiers sont modifiés → Tapez 'F' pour les ajouter au stage")
	}

	if staged != "" {
		suggestions = append(suggestions, "✅ Des fichiers sont en stage → Tapez 'C' pour créer un commit")
	}

	// Suggestions basées sur la branche
	if currentBranch == "main" || currentBranch == "master" {
		if status != "" || staged != "" {
			suggestions = append(suggestions, "⚠️  Vous développez sur la branche principale → Tapez 'B' pour créer une feature branch")
		}
	}

	// Suggestions basées sur la synchronisation
	if ahead != "" && ahead != "0" {
		suggestions = append(suggestions, "📤 Vous avez des commits locaux → Tapez 'R' puis '1' pour pusher")
	}

	if behind != "" && behind != "0" {
		suggestions = append(suggestions, "📥 Des commits sont disponibles sur le remote → Tapez 'R' puis '2' pour puller")
	}

	// Suggestions générales
	if len(suggestions) == 0 {
		if status == "" && staged == "" {
			suggestions = append(suggestions, "🎉 Working directory clean → Bon moment pour créer une nouvelle branche")
			suggestions = append(suggestions, "📊 Tapez '8' pour voir les statistiques du projet")
		}
	}

	// Suggestions de bonnes pratiques
	lastCommitRelative, _ := gm.runGitCommand("log", "-1", "--pretty=format:%ar")
	if strings.Contains(lastCommitRelative, "day") || strings.Contains(lastCommitRelative, "week") || strings.Contains(lastCommitRelative, "month") || strings.Contains(lastCommitRelative, "year") {
		suggestions = append(suggestions, "🕐 Dernier commit ancien → Pensez à faire des commits plus fréquents")
	}

	// Afficher les suggestions
	if len(suggestions) > 0 {
		for i, suggestion := range suggestions {
			if i < 4 { // Limiter à 4 suggestions max
				fmt.Printf("   %s%s%s\n", ColorYellow, suggestion, ColorReset)
			}
		}
	} else {
		fmt.Printf("   %s🎯 Tout semble en ordre! Continuez le bon travail.%s\n", ColorGreen, ColorReset)
	}

	fmt.Println()

	// Actions rapides recommandées
	fmt.Printf("%s%s⚡ ACTIONS RAPIDES DISPONIBLES:%s\n", ColorBold, ColorCyan, ColorReset)
	fmt.Printf("   %sS%s = Actualiser ce statut  %sC%s = Commits  %sF%s = Fichiers  %sB%s = Branches  %sR%s = Remote\n",
		ColorCyan, ColorReset, ColorCyan, ColorReset, ColorCyan, ColorReset, ColorCyan, ColorReset, ColorCyan, ColorReset)
}

func (gm *GitManager) handleBranchManagement() {
	if !gm.isGitRepo() {
		fmt.Printf("%s❌ Ce répertoire n'est pas un dépôt Git!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	for {
		gm.clearScreen()
		fmt.Printf("%s%s🌿 GESTION DES BRANCHES%s\n", ColorBold, ColorGreen, ColorReset)
		fmt.Println(strings.Repeat("═", 30))

		branches, _ := gm.runGitCommand("branch", "-v")
		fmt.Printf("%sBranches locales:%s\n", ColorBlue, ColorReset)
		fmt.Println(branches)

		fmt.Println("\n1. Créer une nouvelle branche")
		fmt.Println("2. Changer de branche")
		fmt.Println("3. Supprimer une branche")
		fmt.Println("4. Renommer une branche")
		fmt.Println("5. Merger une branche")
		fmt.Println("6. Voir les branches remote")
		fmt.Println("0. Retour au menu principal")

		fmt.Printf("\n%sChoisissez une option: %s", ColorYellow, ColorReset)
		choice := gm.getUserInput()

		switch choice {
		case "1":
			gm.createBranch()
		case "2":
			gm.switchBranch()
		case "3":
			gm.deleteBranch()
		case "4":
			gm.renameBranch()
		case "5":
			gm.mergeBranch()
		case "6":
			gm.showRemoteBranches()
		case "0":
			return
		default:
			fmt.Printf("%s❌ Option invalide!%s\n", ColorRed, ColorReset)
			gm.pause()
		}
	}
}

func (gm *GitManager) handleCommitManagement() {
	if !gm.isGitRepo() {
		fmt.Printf("%s❌ Ce répertoire n'est pas un dépôt Git!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	for {
		gm.clearScreen()
		fmt.Printf("%s%s📦 GESTION DES COMMITS%s\n", ColorBold, ColorGreen, ColorReset)
		fmt.Println(strings.Repeat("═", 30))

		fmt.Println("1. Faire un commit")
		fmt.Println("2. Voir l'historique des commits")
		fmt.Println("3. Voir les détails d'un commit")
		fmt.Println("4. Modifier le dernier commit (amend)")
		fmt.Println("5. Reset / Revert")
		fmt.Println("6. Chercher dans les commits")
		fmt.Println("0. Retour au menu principal")

		fmt.Printf("\n%sChoisissez une option: %s", ColorYellow, ColorReset)
		choice := gm.getUserInput()

		switch choice {
		case "1":
			gm.makeCommit()
		case "2":
			gm.showDetailedHistory()
		case "3":
			gm.showCommitDetails()
		case "4":
			gm.amendCommit()
		case "5":
			gm.handleResetRevert()
		case "6":
			gm.searchHistory()
		case "0":
			return
		default:
			fmt.Printf("%s❌ Option invalide!%s\n", ColorRed, ColorReset)
			gm.pause()
		}
	}
}

func (gm *GitManager) handleRemoteManagement() {
	if !gm.isGitRepo() {
		fmt.Printf("%s❌ Ce répertoire n'est pas un dépôt Git!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	for {
		gm.clearScreen()
		fmt.Printf("%s%s🔄 GESTION DES REMOTES%s\n", ColorBold, ColorPurple, ColorReset)
		fmt.Println(strings.Repeat("═", 30))

		remotes, _ := gm.runGitCommand("remote", "-v")
		if remotes != "" {
			fmt.Printf("%sRemotes configurés:%s\n", ColorBlue, ColorReset)
			fmt.Println(remotes)
		} else {
			fmt.Printf("%sAucun remote configuré.%s\n", ColorYellow, ColorReset)
		}

		fmt.Println("\n1. Ajouter un remote")
		fmt.Println("2. Supprimer un remote")
		fmt.Println("3. Renommer un remote")
		fmt.Println("4. Fetch depuis un remote")
		fmt.Println("5. Pull depuis un remote")
		fmt.Println("6. Push vers un remote")
		fmt.Println("0. Retour au menu principal")

		fmt.Printf("\n%sChoisissez une option: %s", ColorYellow, ColorReset)
		choice := gm.getUserInput()

		switch choice {
		case "1":
			gm.addRemote()
		case "2":
			gm.removeRemote()
		case "3":
			gm.renameRemote()
		case "4":
			gm.fetchFromRemote()
		case "5":
			gm.pullFromRemote()
		case "6":
			gm.pushToRemote()
		case "0":
			return
		default:
			fmt.Printf("%s❌ Option invalide!%s\n", ColorRed, ColorReset)
			gm.pause()
		}
	}
}

func (gm *GitManager) handleFileManagement() {
	if !gm.isGitRepo() {
		fmt.Printf("%s❌ Ce répertoire n'est pas un dépôt Git!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	for {
		gm.clearScreen()
		fmt.Printf("%s%s📁 GESTION DES FICHIERS%s\n", ColorBold, ColorGreen, ColorReset)
		fmt.Println(strings.Repeat("═", 30))

		fmt.Println("1. Ajouter des fichiers (add)")
		fmt.Println("2. Retirer des fichiers du staging (reset)")
		fmt.Println("3. Voir les différences (diff)")
		fmt.Println("4. Gérer .gitignore")
		fmt.Println("5. Restaurer des fichiers (checkout)")
		fmt.Println("6. Ne plus tracker des fichiers (rm)")
		fmt.Println("0. Retour au menu principal")

		fmt.Printf("\n%sChoisissez une option: %s", ColorYellow, ColorReset)
		choice := gm.getUserInput()

		switch choice {
		case "1":
			gm.addFiles()
		case "2":
			gm.unstageFiles()
		case "3":
			gm.showDiff()
		case "4":
			gm.manageGitignore()
		case "5":
			gm.restoreFiles()
		case "6":
			gm.untrackFiles()
		case "0":
			return
		default:
			fmt.Printf("%s❌ Option invalide!%s\n", ColorRed, ColorReset)
			gm.pause()
		}
	}
}

func (gm *GitManager) handleTagManagement() {
	if !gm.isGitRepo() {
		fmt.Printf("%s❌ Ce répertoire n'est pas un dépôt Git!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	for {
		gm.clearScreen()
		fmt.Printf("%s%s🏷️  GESTION DES TAGS%s\n", ColorBold, ColorGreen, ColorReset)
		fmt.Println(strings.Repeat("═", 25))

		tags, _ := gm.runGitCommand("tag", "-l")
		if tags != "" {
			fmt.Printf("%s🏷️  Tags existants:%s\n", ColorBlue, ColorReset)
			fmt.Println(tags)
			fmt.Println()
		}

		fmt.Println("1. Créer un tag")
		fmt.Println("2. Créer un tag annoté")
		fmt.Println("3. Supprimer un tag")
		fmt.Println("4. Voir les détails d'un tag")
		fmt.Println("5. Lister tous les tags")
		fmt.Println("0. Retour au menu principal")

		fmt.Printf("\n%sChoisissez une option: %s", ColorYellow, ColorReset)
		choice := gm.getUserInput()

		switch choice {
		case "1":
			gm.createTag()
		case "2":
			gm.createAnnotatedTag()
		case "3":
			gm.deleteTag()
		case "4":
			gm.showTagDetails()
		case "5":
			gm.listTags()
		case "0":
			return
		default:
			fmt.Printf("%s❌ Option invalide!%s\n", ColorRed, ColorReset)
			gm.pause()
		}
	}
}

func (gm *GitManager) handleStashManagement() {
	if !gm.isGitRepo() {
		fmt.Printf("%s❌ Ce répertoire n'est pas un dépôt Git!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	for {
		gm.clearScreen()
		fmt.Printf("%s%s🗂️  GESTION DES STASH%s\n", ColorBold, ColorGreen, ColorReset)
		fmt.Println(strings.Repeat("═", 25))

		stashes, _ := gm.runGitCommand("stash", "list")
		if stashes != "" {
			fmt.Printf("%s🗂️  Stashes existants:%s\n", ColorBlue, ColorReset)
			fmt.Println(stashes)
			fmt.Println()
		}

		fmt.Println("1. Créer un stash")
		fmt.Println("2. Appliquer un stash")
		fmt.Println("3. Voir le contenu d'un stash")
		fmt.Println("4. Supprimer un stash")
		fmt.Println("5. Supprimer tous les stashes")
		fmt.Println("6. Créer une branche depuis un stash")
		fmt.Println("0. Retour au menu principal")

		fmt.Printf("\n%sChoisissez une option: %s", ColorYellow, ColorReset)
		choice := gm.getUserInput()

		switch choice {
		case "1":
			gm.createStash()
		case "2":
			gm.applyStash()
		case "3":
			gm.showStash()
		case "4":
			gm.dropStash()
		case "5":
			gm.clearStashes()
		case "6":
			gm.stashToBranch()
		case "0":
			return
		default:
			fmt.Printf("%s❌ Option invalide!%s\n", ColorRed, ColorReset)
			gm.pause()
		}
	}
}

func (gm *GitManager) handleStatistics() {
	if !gm.isGitRepo() {
		fmt.Printf("%s❌ Ce répertoire n'est pas un dépôt Git!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	for {
		gm.clearScreen()
		fmt.Printf("%s%s📈 STATISTIQUES ET LOGS%s\n", ColorBold, ColorGreen, ColorReset)
		fmt.Println(strings.Repeat("═", 30))

		fmt.Println("1. Statistiques générales")
		fmt.Println("2. Contributeurs et activité")
		fmt.Println("3. Historique détaillé")
		fmt.Println("4. Graphique des branches")
		fmt.Println("5. Statistiques par fichier")
		fmt.Println("6. Blâme d'un fichier")
		fmt.Println("7. Recherche dans l'historique")
		fmt.Println("0. Retour au menu principal")

		fmt.Printf("\n%sChoisissez une option: %s", ColorYellow, ColorReset)
		choice := gm.getUserInput()

		switch choice {
		case "1":
			gm.showGeneralStats()
		case "2":
			gm.showContributorStats()
		case "3":
			gm.showDetailedHistory()
		case "4":
			gm.showBranchGraph()
		case "5":
			gm.showFileStats()
		case "6":
			gm.showBlame()
		case "7":
			gm.searchHistory()
		case "0":
			return
		default:
			fmt.Printf("%s❌ Option invalide!%s\n", ColorRed, ColorReset)
			gm.pause()
		}
	}
}

func (gm *GitManager) handleTools() {
	if !gm.isGitRepo() {
		fmt.Printf("%s❌ Ce répertoire n'est pas un dépôt Git!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	for {
		gm.clearScreen()
		fmt.Printf("%s%s🔧 OUTILS ET CONFIGURATION%s\n", ColorBold, ColorGreen, ColorReset)
		fmt.Println(strings.Repeat("═", 35))

		fmt.Println("1. Configuration Git")
		fmt.Println("2. Nettoyage du dépôt")
		fmt.Println("3. Vérification du dépôt")
		fmt.Println("4. Hooks Git")
		fmt.Println("5. Aliases Git")
		fmt.Println("6. Sauvegarde/Archive")
		fmt.Println("0. Retour au menu principal")

		fmt.Printf("\n%sChoisissez une option: %s", ColorYellow, ColorReset)
		choice := gm.getUserInput()

		switch choice {
		case "1":
			gm.manageConfig()
		case "2":
			gm.cleanRepo()
		case "3":
			gm.checkRepo()
		case "4":
			gm.manageHooks()
		case "5":
			gm.manageAliases()
		case "6":
			gm.archiveRepo()
		case "0":
			return
		default:
			fmt.Printf("%s❌ Option invalide!%s\n", ColorRed, ColorReset)
			gm.pause()
		}
	}
}

// Branch Management
func (gm *GitManager) createBranch() {
	fmt.Printf("%sNom de la nouvelle branche: %s", ColorYellow, ColorReset)
	branchName := gm.getUserInput()

	if branchName == "" {
		fmt.Printf("%s❌ Nom de branche invalide!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	output, err := gm.runGitCommand("checkout", "-b", branchName)
	if err != nil {
		fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
	} else {
		fmt.Printf("%s✅ Branche '%s' créée et activée!%s\n", ColorGreen, branchName, ColorReset)
	}
	gm.pause()
}

func (gm *GitManager) switchBranch() {
	branches, _ := gm.runGitCommand("branch")
	fmt.Printf("%sBranches disponibles:%s\n", ColorBlue, ColorReset)
	fmt.Println(branches)

	fmt.Printf("\n%sNom de la branche: %s", ColorYellow, ColorReset)
	branchName := gm.getUserInput()

	if branchName == "" {
		fmt.Printf("%s❌ Nom de branche invalide!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	output, err := gm.runGitCommand("checkout", branchName)
	if err != nil {
		fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
	} else {
		fmt.Printf("%s✅ Branche '%s' activée!%s\n", ColorGreen, branchName, ColorReset)
	}
	gm.pause()
}

func (gm *GitManager) deleteBranch() {
	branches, _ := gm.runGitCommand("branch")
	fmt.Printf("%sBranches disponibles:%s\n", ColorBlue, ColorReset)
	fmt.Println(branches)

	fmt.Printf("\n%sNom de la branche à supprimer: %s", ColorYellow, ColorReset)
	branchName := gm.getUserInput()

	if branchName == "" {
		fmt.Printf("%s❌ Nom de branche invalide!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	fmt.Printf("%s⚠️  Êtes-vous sûr de vouloir supprimer '%s'? (y/N): %s", ColorRed, branchName, ColorReset)
	confirm := gm.getUserInput()

	if strings.ToLower(confirm) == "y" {
		output, err := gm.runGitCommand("branch", "-d", branchName)
		if err != nil {
			fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
			fmt.Printf("%s💡 Utilisez 'git branch -D %s' pour forcer la suppression%s\n", ColorYellow, branchName, ColorReset)
		} else {
			fmt.Printf("%s✅ Branche '%s' supprimée!%s\n", ColorGreen, branchName, ColorReset)
		}
	}
	gm.pause()
}

func (gm *GitManager) renameBranch() {
	currentBranch := gm.getCurrentBranch()
	fmt.Printf("%sBranche actuelle: %s%s%s\n", ColorBlue, ColorCyan, currentBranch, ColorReset)

	fmt.Printf("%sNouveau nom: %s", ColorYellow, ColorReset)
	newName := gm.getUserInput()

	if newName == "" {
		fmt.Printf("%s❌ Nom invalide!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	output, err := gm.runGitCommand("branch", "-m", newName)
	if err != nil {
		fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
	} else {
		fmt.Printf("%s✅ Branche renommée de '%s' à '%s'!%s\n", ColorGreen, currentBranch, newName, ColorReset)
	}
	gm.pause()
}

func (gm *GitManager) mergeBranch() {
	currentBranch := gm.getCurrentBranch()
	branches, _ := gm.runGitCommand("branch")

	fmt.Printf("%sBranche actuelle: %s%s%s\n", ColorBlue, ColorCyan, currentBranch, ColorReset)
	fmt.Printf("%sBranches disponibles:%s\n", ColorBlue, ColorReset)
	fmt.Println(branches)

	fmt.Printf("\n%sBranche à merger dans '%s': %s", ColorYellow, currentBranch, ColorReset)
	branchName := gm.getUserInput()

	if branchName == "" {
		fmt.Printf("%s❌ Nom de branche invalide!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	output, err := gm.runGitCommand("merge", branchName)
	if err != nil {
		fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
	} else {
		fmt.Printf("%s✅ Branche '%s' mergée dans '%s'!%s\n", ColorGreen, branchName, currentBranch, ColorReset)
	}
	gm.pause()
}

func (gm *GitManager) showRemoteBranches() {
	output, err := gm.runGitCommand("branch", "-r")
	if err != nil {
		fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
	} else {
		fmt.Printf("%s🌐 Branches remote:%s\n", ColorCyan, ColorReset)
		fmt.Println(output)
	}
	gm.pause()
}

// Commit Management
func (gm *GitManager) makeCommit() {
	staged, _ := gm.runGitCommand("diff", "--cached", "--name-only")
	if staged == "" {
		fmt.Printf("%sAucun fichier en stage. Voulez-vous ajouter des fichiers? (y/N): %s", ColorYellow, ColorReset)
		choice := gm.getUserInput()
		if strings.ToLower(choice) == "y" {
			gm.addFiles()
			staged, _ = gm.runGitCommand("diff", "--cached", "--name-only")
			if staged == "" {
				fmt.Printf("%s❌ Aucun fichier en stage après ajout. Annulation du commit.%s\n", ColorRed, ColorReset)
				gm.pause()
				return
			}
		} else {
			gm.pause()
			return
		}
	}

	fmt.Printf("%sMessage de commit: %s", ColorYellow, ColorReset)
	message := gm.getUserInput()

	if message == "" {
		fmt.Printf("%s❌ Message de commit requis!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	output, err := gm.runGitCommand("commit", "-m", message)
	if err != nil {
		fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
	} else {
		fmt.Printf("%s✅ Commit créé!%s\n", ColorGreen, ColorReset)
		fmt.Println(output)
	}
	gm.pause()
}

func (gm *GitManager) showCommitDetails() {
	fmt.Printf("%sHash du commit (vide pour le dernier): %s", ColorYellow, ColorReset)
	commitHash := gm.getUserInput()

	var output string
	var err error

	// Ajout de l'option --color=always pour forcer la coloration
	if commitHash == "" {
		output, err = gm.runGitCommand("show", "--color=always", "HEAD")
	} else {
		output, err = gm.runGitCommand("show", "--color=always", commitHash)
	}

	if err != nil {
		fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
	} else {
		fmt.Println(output)
	}
	gm.pause()
}

func (gm *GitManager) amendCommit() {
	fmt.Printf("%sVoulez-vous modifier le message du dernier commit? (y/N): %s", ColorYellow, ColorReset)
	choice := gm.getUserInput()
	var output string
	var err error

	if strings.ToLower(choice) == "y" {
		fmt.Printf("%sNouveau message de commit: %s", ColorYellow, ColorReset)
		newMessage := gm.getUserInput()
		if newMessage == "" {
			fmt.Printf("%s❌ Message de commit requis!%s\n", ColorRed, ColorReset)
			gm.pause()
			return
		}
		output, err = gm.runGitCommand("commit", "--amend", "-m", newMessage)
	} else {
		output, err = gm.runGitCommand("commit", "--amend", "--no-edit")
	}

	if err != nil {
		fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
	} else {
		fmt.Printf("%s✅ Dernier commit modifié!%s\n", ColorGreen, ColorReset)
	}
	gm.pause()
}

func (gm *GitManager) handleResetRevert() {
	for {
		gm.clearScreen()
		fmt.Printf("%s%s🔄 RESET / REVERT%s\n", ColorBold, ColorRed, ColorReset)
		fmt.Println(strings.Repeat("═", 30))
		fmt.Println("1. Reset (déplacer HEAD)")
		fmt.Println("2. Revert (créer un commit d'annulation)")
		fmt.Println("0. Retour")

		fmt.Printf("\n%sChoisissez une option: %s", ColorYellow, ColorReset)
		choice := gm.getUserInput()

		switch choice {
		case "1":
			gm.handleReset()
		case "2":
			gm.handleRevert()
		case "0":
			return
		default:
			fmt.Printf("%s❌ Option invalide!%s\n", ColorRed, ColorReset)
			gm.pause()
		}
	}
}

func (gm *GitManager) handleReset() {
	fmt.Printf("%sCommit cible (ex: HEAD~1, hash): %s", ColorYellow, ColorReset)
	target := gm.getUserInput()
	if target == "" {
		fmt.Printf("%s❌ Cible requise!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	fmt.Println("Types de reset:")
	fmt.Println("1. --soft (ne touche pas à l'index ni à l'arbre de travail)")
	fmt.Println("2. --mixed (défaut, reset l'index mais pas l'arbre de travail)")
	fmt.Println("3. --hard (ATTENTION: perd les modifications locales)")

	fmt.Printf("\n%sChoisissez un type de reset (défaut 2): %s", ColorYellow, ColorReset)
	choice := gm.getUserInput()

	var resetType string
	switch choice {
	case "1":
		resetType = "--soft"
	case "3":
		fmt.Printf("%s⚠️  Êtes-vous sûr de vouloir faire un reset --hard? (y/N): %s", ColorRed, ColorReset)
		confirm := gm.getUserInput()
		if strings.ToLower(confirm) != "y" {
			gm.pause()
			return
		}
		resetType = "--hard"
	default:
		resetType = "--mixed"
	}

	output, err := gm.runGitCommand("reset", resetType, target)
	if err != nil {
		fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
	} else {
		fmt.Printf("%s✅ Reset effectué!%s\n", ColorGreen, ColorReset)
		fmt.Println(output)
	}
	gm.pause()
}

func (gm *GitManager) handleRevert() {
	fmt.Printf("%sCommit à annuler (revert): %s", ColorYellow, ColorReset)
	target := gm.getUserInput()
	if target == "" {
		fmt.Printf("%s❌ Cible requise!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	output, err := gm.runGitCommand("revert", target)
	if err != nil {
		fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
	} else {
		fmt.Printf("%s✅ Revert effectué!%s\n", ColorGreen, ColorReset)
		fmt.Println(output)
	}
	gm.pause()
}

// Remote Management
func (gm *GitManager) addRemote() {
	fmt.Printf("%sNom du remote (ex: origin): %s", ColorYellow, ColorReset)
	name := gm.getUserInput()
	fmt.Printf("%sURL du remote: %s", ColorYellow, ColorReset)
	url := gm.getUserInput()

	if name == "" || url == "" {
		fmt.Printf("%s❌ Nom et URL sont requis!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	output, err := gm.runGitCommand("remote", "add", name, url)
	if err != nil {
		fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
	} else {
		fmt.Printf("%s✅ Remote '%s' ajouté!%s\n", ColorGreen, name, ColorReset)
	}
	gm.pause()
}

func (gm *GitManager) removeRemote() {
	fmt.Printf("%sNom du remote à supprimer: %s", ColorYellow, ColorReset)
	name := gm.getUserInput()
	if name == "" {
		gm.pause()
		return
	}

	output, err := gm.runGitCommand("remote", "remove", name)
	if err != nil {
		fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
	} else {
		fmt.Printf("%s✅ Remote '%s' supprimé!%s\n", ColorGreen, name, ColorReset)
	}
	gm.pause()
}

func (gm *GitManager) renameRemote() {
	fmt.Printf("%sAncien nom du remote: %s", ColorYellow, ColorReset)
	oldName := gm.getUserInput()
	fmt.Printf("%sNouveau nom du remote: %s", ColorYellow, ColorReset)
	newName := gm.getUserInput()

	if oldName == "" || newName == "" {
		fmt.Printf("%s❌ Les deux noms sont requis!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	output, err := gm.runGitCommand("remote", "rename", oldName, newName)
	if err != nil {
		fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
	} else {
		fmt.Printf("%s✅ Remote renommé de '%s' à '%s'!%s\n", ColorGreen, oldName, newName, ColorReset)
	}
	gm.pause()
}

func (gm *GitManager) fetchFromRemote() {
	fmt.Printf("%sNom du remote (laisser vide pour 'all'): %s", ColorYellow, ColorReset)
	remote := gm.getUserInput()

	var output string
	var err error
	if remote == "" {
		output, err = gm.runGitCommand("fetch", "--all", "--prune")
	} else {
		output, err = gm.runGitCommand("fetch", remote, "--prune")
	}

	if err != nil {
		fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
	} else {
		fmt.Printf("%s✅ Fetch terminé!%s\n", ColorGreen, ColorReset)
		fmt.Println(output)
	}
	gm.pause()
}

func (gm *GitManager) pullFromRemote() {
	fmt.Printf("%sRemote (défaut 'origin'): %s", ColorYellow, ColorReset)
	remote := gm.getUserInput()
	if remote == "" {
		remote = "origin"
	}

	currentBranch := gm.getCurrentBranch()
	fmt.Printf("%sBranche (défaut '%s'): %s", ColorYellow, currentBranch, ColorReset)
	branch := gm.getUserInput()
	if branch == "" {
		branch = currentBranch
	}

	output, err := gm.runGitCommand("pull", remote, branch)
	if err != nil {
		fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
	} else {
		fmt.Printf("%s✅ Pull terminé!%s\n", ColorGreen, ColorReset)
		fmt.Println(output)
	}
	gm.pause()
}

func (gm *GitManager) pushToRemote() {
	fmt.Printf("%sRemote (défaut 'origin'): %s", ColorYellow, ColorReset)
	remote := gm.getUserInput()
	if remote == "" {
		remote = "origin"
	}

	currentBranch := gm.getCurrentBranch()
	fmt.Printf("%sBranche (défaut '%s'): %s", ColorYellow, currentBranch, ColorReset)
	branch := gm.getUserInput()
	if branch == "" {
		branch = currentBranch
	}

	fmt.Printf("%sForcer le push? (y/N): %s", ColorRed, ColorReset)
	force := gm.getUserInput()

	args := []string{"push", remote, branch}
	if strings.ToLower(force) == "y" {
		args = append(args, "--force")
	}

	output, err := gm.runGitCommand(args...)
	if err != nil {
		fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
	} else {
		fmt.Printf("%s✅ Push terminé!%s\n", ColorGreen, ColorReset)
		fmt.Println(output)
	}
	gm.pause()
}

// File Management
func (gm *GitManager) printFileStatus(line string) {
	if len(line) < 3 {
		return
	}

	status := line[:2]
	file := line[3:]

	switch status {
	case "M ":
		fmt.Printf("%s  M  %s%s (modifié, en stage)\n", ColorGreen, file, ColorReset)
	case " M":
		fmt.Printf("%s  M  %s%s (modifié)\n", ColorYellow, file, ColorReset)
	case "A ":
		fmt.Printf("%s  A  %s%s (ajouté)\n", ColorGreen, file, ColorReset)
	case "D ":
		fmt.Printf("%s  D  %s%s (supprimé)\n", ColorRed, file, ColorReset)
	case "??":
		fmt.Printf("%s  ?  %s%s (non suivi)\n", ColorRed, file, ColorReset)
	case "R ":
		fmt.Printf("%s  R  %s%s (renommé)\n", ColorCyan, file, ColorReset)
	default:
		fmt.Printf("%s  %s  %s%s\n", ColorWhite, status, file, ColorReset)
	}
}

func (gm *GitManager) addFiles() {
	status := gm.getGitStatus()
	if status == "" {
		fmt.Printf("%s✅ Aucun fichier à ajouter!%s\n", ColorGreen, ColorReset)
		gm.pause()
		return
	}

	fmt.Printf("%s📁 Fichiers disponibles:%s\n", ColorYellow, ColorReset)
	lines := strings.Split(status, "\n")
	for _, line := range lines {
		if line != "" {
			gm.printFileStatus(line)
		}
	}

	fmt.Println("\n1. Ajouter tous les fichiers")
	fmt.Println("2. Ajouter des fichiers spécifiques")
	fmt.Println("0. Annuler")

	fmt.Printf("\n%sChoisissez une option: %s", ColorYellow, ColorReset)
	choice := gm.getUserInput()

	switch choice {
	case "1":
		output, err := gm.runGitCommand("add", ".")
		if err != nil {
			fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
		} else {
			fmt.Printf("%s✅ Tous les fichiers ajoutés!%s\n", ColorGreen, ColorReset)
		}
	case "2":
		fmt.Printf("%sFichiers à ajouter (séparés par des espaces): %s", ColorYellow, ColorReset)
		files := gm.getUserInput()
		if files != "" {
			fileList := strings.Fields(files)
			for _, file := range fileList {
				output, err := gm.runGitCommand("add", file)
				if err != nil {
					fmt.Printf("%s❌ Erreur avec '%s': %s%s\n", ColorRed, file, output, ColorReset)
				} else {
					fmt.Printf("%s✅ '%s' ajouté!%s\n", ColorGreen, file, ColorReset)
				}
			}
		}
	}
	gm.pause()
}

func (gm *GitManager) unstageFiles() {
	staged, _ := gm.runGitCommand("diff", "--cached", "--name-only")
	if staged == "" {
		fmt.Printf("%s✅ Aucun fichier en stage!%s\n", ColorGreen, ColorReset)
		gm.pause()
		return
	}

	fmt.Printf("%s📁 Fichiers en stage:%s\n", ColorYellow, ColorReset)
	fmt.Println(staged)

	fmt.Println("\n1. Retirer tous les fichiers du staging")
	fmt.Println("2. Retirer des fichiers spécifiques")
	fmt.Println("0. Annuler")

	fmt.Printf("\n%sChoisissez une option: %s", ColorYellow, ColorReset)
	choice := gm.getUserInput()

	switch choice {
	case "1":
		output, err := gm.runGitCommand("reset", "HEAD")
		if err != nil {
			fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
		} else {
			fmt.Printf("%s✅ Tous les fichiers retirés du staging!%s\n", ColorGreen, ColorReset)
		}
	case "2":
		fmt.Printf("%sFichiers à retirer (séparés par des espaces): %s", ColorYellow, ColorReset)
		files := gm.getUserInput()
		if files != "" {
			fileList := strings.Fields(files)
			for _, file := range fileList {
				output, err := gm.runGitCommand("reset", "HEAD", "--", file)
				if err != nil {
					fmt.Printf("%s❌ Erreur avec '%s': %s%s\n", ColorRed, file, output, ColorReset)
				} else {
					fmt.Printf("%s✅ '%s' retiré du staging!%s\n", ColorGreen, file, ColorReset)
				}
			}
		}
	}
	gm.pause()
}

func (gm *GitManager) showDiff() {
	for {
		gm.clearScreen()
		fmt.Println("1. Voir les changements non stagés")
		fmt.Println("2. Voir les changements stagés")
		fmt.Println("3. Voir les changements d'un fichier spécifique")
		fmt.Println("4. Comparer deux commits")
		fmt.Println("0. Retour")

		fmt.Printf("\n%sChoisissez une option: %s", ColorYellow, ColorReset)
		choice := gm.getUserInput()

		switch choice {
		case "1":
			// Ajout de l'option --color=always pour forcer la coloration
			output, _ := gm.runGitCommand("diff", "--color=always")
			if output == "" {
				fmt.Printf("%s✅ Aucun changement non stagé!%s\n", ColorGreen, ColorReset)
			} else {
				fmt.Printf("%s📊 Changements non stagés:%s\n", ColorBlue, ColorReset)
				fmt.Println(output)
			}
			gm.pause()
		case "2":
			// Ajout de l'option --color=always pour forcer la coloration
			output, _ := gm.runGitCommand("diff", "--cached", "--color=always")
			if output == "" {
				fmt.Printf("%s✅ Aucun changement stagé!%s\n", ColorGreen, ColorReset)
			} else {
				fmt.Printf("%s📊 Changements stagés:%s\n", ColorBlue, ColorReset)
				fmt.Println(output)
			}
			gm.pause()
		case "3":
			fmt.Printf("%sNom du fichier: %s", ColorYellow, ColorReset)
			filename := gm.getUserInput()
			if filename != "" {
				// Ajout de l'option --color=always pour forcer la coloration
				output, err := gm.runGitCommand("diff", "--color=always", "--", filename)
				if err != nil {
					fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
				} else {
					fmt.Printf("%s📊 Changements dans '%s':%s\n", ColorBlue, filename, ColorReset)
					fmt.Println(output)
				}
			}
			gm.pause()
		case "4":
			fmt.Printf("%sPremier commit: %s", ColorYellow, ColorReset)
			commit1 := gm.getUserInput()
			fmt.Printf("%sDeuxième commit: %s", ColorYellow, ColorReset)
			commit2 := gm.getUserInput()

			if commit1 != "" && commit2 != "" {
				// Ajout de l'option --color=always pour forcer la coloration
				output, err := gm.runGitCommand("diff", "--color=always", commit1, commit2)
				if err != nil {
					fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
				} else {
					fmt.Printf("%s📊 Différences entre '%s' et '%s':%s\n", ColorBlue, commit1, commit2, ColorReset)
					fmt.Println(output)
				}
			}
			gm.pause()
		case "0":
			return
		default:
			fmt.Printf("%s❌ Option invalide!%s\n", ColorRed, ColorReset)
			gm.pause()
		}
	}
}

func (gm *GitManager) manageGitignore() {
	gitignorePath := filepath.Join(gm.currentPath, ".gitignore")

	for {
		gm.clearScreen()
		fmt.Println("1. Voir le contenu de .gitignore")
		fmt.Println("2. Ajouter des patterns à .gitignore")
		fmt.Println("3. Créer un .gitignore basique")
		fmt.Println("0. Retour")

		fmt.Printf("\n%sChoisissez une option: %s", ColorYellow, ColorReset)
		choice := gm.getUserInput()

		switch choice {
		case "1":
			content, err := os.ReadFile(gitignorePath)
			if err != nil {
				fmt.Printf("%s❌ Fichier .gitignore introuvable ou erreur de lecture%s\n", ColorRed, ColorReset)
			} else {
				fmt.Printf("%s📄 Contenu de .gitignore:%s\n", ColorBlue, ColorReset)
				fmt.Println(string(content))
			}
			gm.pause()
		case "2":
			fmt.Printf("%sPatterns à ajouter (un par ligne, ligne vide pour terminer):%s\n", ColorYellow, ColorReset)
			var patterns []string
			for {
				fmt.Print("> ")
				line := gm.getUserInput()
				if line == "" {
					break
				}
				patterns = append(patterns, line)
			}

			if len(patterns) > 0 {
				file, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					fmt.Printf("%s❌ Erreur lors de l'ouverture du fichier: %s%s\n", ColorRed, err.Error(), ColorReset)
				} else {
					defer file.Close()
					for _, pattern := range patterns {
						_, _ = file.WriteString(pattern + "\n")
					}
					fmt.Printf("%s✅ Patterns ajoutés à .gitignore!%s\n", ColorGreen, ColorReset)
				}
			}
			gm.pause()
		case "3":
			basicGitignore := `# Fichiers de compilation
*.o
*.so
*.dylib
*.exe

# Fichiers de debug
*.dSYM/
*.pdb

# Fichiers temporaires
*.tmp
*.temp
*~
.DS_Store
Thumbs.db

# Répertoires de dépendances
node_modules/
vendor/

# Fichiers de configuration locale
.env
.env.local
config.local.*

# Logs
*.log
logs/

# Fichiers d'IDE
.vscode/
.idea/
*.swp
*.swo
`
			err := os.WriteFile(gitignorePath, []byte(basicGitignore), 0644)
			if err != nil {
				fmt.Printf("%s❌ Erreur lors de la création: %s%s\n", ColorRed, err.Error(), ColorReset)
			} else {
				fmt.Printf("%s✅ .gitignore basique créé!%s\n", ColorGreen, ColorReset)
			}
			gm.pause()
		case "0":
			return
		default:
			fmt.Printf("%s❌ Option invalide!%s\n", ColorRed, ColorReset)
			gm.pause()
		}
	}
}

func (gm *GitManager) restoreFiles() {
	status := gm.getGitStatus()
	if status == "" {
		fmt.Printf("%s✅ Aucun fichier à restaurer!%s\n", ColorGreen, ColorReset)
		gm.pause()
		return
	}

	fmt.Printf("%s📁 Fichiers modifiés:%s\n", ColorYellow, ColorReset)
	lines := strings.Split(status, "\n")
	for _, line := range lines {
		if line != "" {
			gm.printFileStatus(line)
		}
	}

	fmt.Println("\n1. Restaurer tous les fichiers modifiés")
	fmt.Println("2. Restaurer des fichiers spécifiques")
	fmt.Println("0. Annuler")

	fmt.Printf("\n%s⚠️  ATTENTION: Cette action va perdre les modifications non commitées!%s\n", ColorRed, ColorReset)
	fmt.Printf("%sChoisissez une option: %s", ColorYellow, ColorReset)
	choice := gm.getUserInput()

	switch choice {
	case "1":
		fmt.Printf("%sÊtes-vous sûr? (y/N): %s", ColorRed, ColorReset)
		confirm := gm.getUserInput()
		if strings.ToLower(confirm) == "y" {
			output, err := gm.runGitCommand("checkout", "--", ".")
			if err != nil {
				fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
			} else {
				fmt.Printf("%s✅ Tous les fichiers restaurés!%s\n", ColorGreen, ColorReset)
			}
		}
	case "2":
		fmt.Printf("%sFichiers à restaurer (séparés par des espaces): %s", ColorYellow, ColorReset)
		files := gm.getUserInput()
		if files != "" {
			fmt.Printf("%sÊtes-vous sûr? (y/N): %s", ColorRed, ColorReset)
			confirm := gm.getUserInput()
			if strings.ToLower(confirm) == "y" {
				fileList := strings.Fields(files)
				for _, file := range fileList {
					output, err := gm.runGitCommand("checkout", "--", file)
					if err != nil {
						fmt.Printf("%s❌ Erreur avec '%s': %s%s\n", ColorRed, file, output, ColorReset)
					} else {
						fmt.Printf("%s✅ '%s' restauré!%s\n", ColorGreen, file, ColorReset)
					}
				}
			}
		}
	}
	gm.pause()
}

func (gm *GitManager) untrackFiles() {
	fmt.Printf("%sFichiers à ne plus tracker (séparés par des espaces): %s", ColorYellow, ColorReset)
	files := gm.getUserInput()

	if files != "" {
		fmt.Printf("%sGarder les fichiers localement? (Y/n): %s", ColorYellow, ColorReset)
		keep := gm.getUserInput()

		fileList := strings.Fields(files)
		for _, file := range fileList {
			var output string
			var err error

			if strings.ToLower(keep) == "n" {
				output, err = gm.runGitCommand("rm", file)
			} else {
				output, err = gm.runGitCommand("rm", "--cached", file)
			}

			if err != nil {
				fmt.Printf("%s❌ Erreur avec '%s': %s%s\n", ColorRed, file, output, ColorReset)
			} else {
				fmt.Printf("%s✅ '%s' retiré du tracking!%s\n", ColorGreen, file, ColorReset)
			}
		}
	}
	gm.pause()
}

// Tag Management
func (gm *GitManager) createTag() {
	fmt.Printf("%sNom du tag: %s", ColorYellow, ColorReset)
	tagName := gm.getUserInput()

	if tagName == "" {
		fmt.Printf("%s❌ Nom de tag requis!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	fmt.Printf("%sCommit à tagger (vide pour HEAD): %s", ColorYellow, ColorReset)
	commit := gm.getUserInput()

	var output string
	var err error

	if commit == "" {
		output, err = gm.runGitCommand("tag", tagName)
	} else {
		output, err = gm.runGitCommand("tag", tagName, commit)
	}

	if err != nil {
		fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
	} else {
		fmt.Printf("%s✅ Tag '%s' créé!%s\n", ColorGreen, tagName, ColorReset)
	}
	gm.pause()
}

func (gm *GitManager) createAnnotatedTag() {
	fmt.Printf("%sNom du tag: %s", ColorYellow, ColorReset)
	tagName := gm.getUserInput()

	if tagName == "" {
		fmt.Printf("%s❌ Nom de tag requis!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	fmt.Printf("%sMessage du tag: %s", ColorYellow, ColorReset)
	message := gm.getUserInput()

	if message == "" {
		fmt.Printf("%s❌ Message requis pour un tag annoté!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	fmt.Printf("%sCommit à tagger (vide pour HEAD): %s", ColorYellow, ColorReset)
	commit := gm.getUserInput()

	var output string
	var err error

	if commit == "" {
		output, err = gm.runGitCommand("tag", "-a", tagName, "-m", message)
	} else {
		output, err = gm.runGitCommand("tag", "-a", tagName, "-m", message, commit)
	}

	if err != nil {
		fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
	} else {
		fmt.Printf("%s✅ Tag annoté '%s' créé!%s\n", ColorGreen, tagName, ColorReset)
	}
	gm.pause()
}

func (gm *GitManager) deleteTag() {
	tags, _ := gm.runGitCommand("tag", "-l")
	if tags == "" {
		fmt.Printf("%s❌ Aucun tag à supprimer!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	fmt.Printf("%s🏷️  Tags disponibles:%s\n", ColorBlue, ColorReset)
	fmt.Println(tags)

	fmt.Printf("\n%sNom du tag à supprimer: %s", ColorYellow, ColorReset)
	tagName := gm.getUserInput()

	if tagName == "" {
		gm.pause()
		return
	}

	fmt.Printf("%s⚠️  Êtes-vous sûr de vouloir supprimer le tag '%s'? (y/N): %s", ColorRed, tagName, ColorReset)
	confirm := gm.getUserInput()

	if strings.ToLower(confirm) == "y" {
		output, err := gm.runGitCommand("tag", "-d", tagName)
		if err != nil {
			fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
		} else {
			fmt.Printf("%s✅ Tag '%s' supprimé!%s\n", ColorGreen, tagName, ColorReset)
		}
	}
	gm.pause()
}

func (gm *GitManager) showTagDetails() {
	fmt.Printf("%sNom du tag: %s", ColorYellow, ColorReset)
	tagName := gm.getUserInput()

	if tagName == "" {
		gm.pause()
		return
	}

	output, err := gm.runGitCommand("show", tagName)
	if err != nil {
		fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
	} else {
		fmt.Printf("%s🏷️  Détails du tag '%s':%s\n", ColorBlue, tagName, ColorReset)
		fmt.Println(output)
	}
	gm.pause()
}

func (gm *GitManager) listTags() {
	output, err := gm.runGitCommand("tag", "-l", "--sort=-version:refname")
	if err != nil {
		fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
	} else {
		fmt.Printf("%s🏷️  Tous les tags:%s\n", ColorBlue, ColorReset)
		if output == "" {
			fmt.Printf("%s❌ Aucun tag trouvé%s\n", ColorRed, ColorReset)
		} else {
			fmt.Println(output)
		}
	}
	gm.pause()
}

// Stash Management
func (gm *GitManager) createStash() {
	status := gm.getGitStatus()
	if status == "" {
		fmt.Printf("%s✅ Aucun changement à stasher!%s\n", ColorGreen, ColorReset)
		gm.pause()
		return
	}

	fmt.Printf("%sMessage pour le stash (optionnel): %s", ColorYellow, ColorReset)
	message := gm.getUserInput()

	fmt.Println("\n1. Stash normal (fichiers trackés)")
	fmt.Println("2. Stash avec fichiers non trackés (-u)")
	fmt.Println("3. Stash avec tous les fichiers (-a)")
	fmt.Println("0. Annuler")

	fmt.Printf("\n%sChoisissez une option: %s", ColorYellow, ColorReset)
	choice := gm.getUserInput()

	var cmd []string
	if message != "" {
		cmd = []string{"stash", "push", "-m", message}
	} else {
		cmd = []string{"stash"}
	}

	switch choice {
	case "1":
		// default command
	case "2":
		cmd = append(cmd, "-u")
	case "3":
		cmd = append(cmd, "-a")
	case "0":
		return
	default:
		fmt.Printf("%s❌ Option invalide!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	output, err := gm.runGitCommand(cmd...)
	if err != nil {
		fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
	} else {
		fmt.Printf("%s✅ Stash créé!%s\n", ColorGreen, ColorReset)
	}
	gm.pause()
}

func (gm *GitManager) applyStash() {
	stashes, _ := gm.runGitCommand("stash", "list")
	if stashes == "" {
		fmt.Printf("%s❌ Aucun stash disponible!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	fmt.Printf("%s🗂️  Stashes disponibles:%s\n", ColorBlue, ColorReset)
	fmt.Println(stashes)

	fmt.Printf("\n%sIndex du stash à appliquer (0 pour le plus récent): %s", ColorYellow, ColorReset)
	indexStr := gm.getUserInput()

	index := 0
	if indexStr != "" {
		if i, err := strconv.Atoi(indexStr); err == nil {
			index = i
		}
	}

	fmt.Println("\n1. Appliquer et garder le stash (apply)")
	fmt.Println("2. Appliquer et supprimer le stash (pop)")
	fmt.Println("0. Annuler")

	fmt.Printf("\n%sChoisissez une option: %s", ColorYellow, ColorReset)
	choice := gm.getUserInput()

	var output string
	var err error
	stashRef := fmt.Sprintf("stash@{%d}", index)

	switch choice {
	case "1":
		output, err = gm.runGitCommand("stash", "apply", stashRef)
	case "2":
		output, err = gm.runGitCommand("stash", "pop", stashRef)
	case "0":
		return
	default:
		fmt.Printf("%s❌ Option invalide!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	if err != nil {
		fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
	} else {
		fmt.Printf("%s✅ Stash appliqué!%s\n", ColorGreen, ColorReset)
	}
	gm.pause()
}

func (gm *GitManager) showStash() {
	stashes, _ := gm.runGitCommand("stash", "list")
	if stashes == "" {
		fmt.Printf("%s❌ Aucun stash disponible!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	fmt.Printf("%s🗂️  Stashes disponibles:%s\n", ColorBlue, ColorReset)
	fmt.Println(stashes)

	fmt.Printf("\n%sIndex du stash à voir (0 pour le plus récent): %s", ColorYellow, ColorReset)
	indexStr := gm.getUserInput()

	index := 0
	if indexStr != "" {
		if i, err := strconv.Atoi(indexStr); err == nil {
			index = i
		}
	}

	stashRef := fmt.Sprintf("stash@{%d}", index)
	output, err := gm.runGitCommand("stash", "show", "-p", stashRef)
	if err != nil {
		fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
	} else {
		fmt.Printf("%s🗂️  Contenu du stash %s:%s\n", ColorBlue, stashRef, ColorReset)
		fmt.Println(output)
	}
	gm.pause()
}

func (gm *GitManager) dropStash() {
	stashes, _ := gm.runGitCommand("stash", "list")
	if stashes == "" {
		fmt.Printf("%s❌ Aucun stash disponible!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	fmt.Printf("%s🗂️  Stashes disponibles:%s\n", ColorBlue, ColorReset)
	fmt.Println(stashes)

	fmt.Printf("\n%sIndex du stash à supprimer (0 pour le plus récent): %s", ColorYellow, ColorReset)
	indexStr := gm.getUserInput()

	index := 0
	if indexStr != "" {
		if i, err := strconv.Atoi(indexStr); err == nil {
			index = i
		}
	}

	stashRef := fmt.Sprintf("stash@{%d}", index)

	fmt.Printf("%s⚠️  Êtes-vous sûr de vouloir supprimer le stash %s? (y/N): %s", ColorRed, stashRef, ColorReset)
	confirm := gm.getUserInput()

	if strings.ToLower(confirm) == "y" {
		output, err := gm.runGitCommand("stash", "drop", stashRef)
		if err != nil {
			fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
		} else {
			fmt.Printf("%s✅ Stash %s supprimé!%s\n", ColorGreen, stashRef, ColorReset)
		}
	}
	gm.pause()
}

func (gm *GitManager) clearStashes() {
	stashes, _ := gm.runGitCommand("stash", "list")
	if stashes == "" {
		fmt.Printf("%s❌ Aucun stash à supprimer!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	fmt.Printf("%s🗂️  Stashes actuels:%s\n", ColorBlue, ColorReset)
	fmt.Println(stashes)

	fmt.Printf("\n%s⚠️  Êtes-vous sûr de vouloir supprimer TOUS les stashes? (y/N): %s", ColorRed, ColorReset)
	confirm := gm.getUserInput()

	if strings.ToLower(confirm) == "y" {
		output, err := gm.runGitCommand("stash", "clear")
		if err != nil {
			fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
		} else {
			fmt.Printf("%s✅ Tous les stashes supprimés!%s\n", ColorGreen, ColorReset)
		}
	}
	gm.pause()
}

func (gm *GitManager) stashToBranch() {
	stashes, _ := gm.runGitCommand("stash", "list")
	if stashes == "" {
		fmt.Printf("%s❌ Aucun stash disponible!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	fmt.Printf("%s🗂️  Stashes disponibles:%s\n", ColorBlue, ColorReset)
	fmt.Println(stashes)

	fmt.Printf("\n%sIndex du stash (0 pour le plus récent): %s", ColorYellow, ColorReset)
	indexStr := gm.getUserInput()

	index := 0
	if indexStr != "" {
		if i, err := strconv.Atoi(indexStr); err == nil {
			index = i
		}
	}

	fmt.Printf("%sNom de la nouvelle branche: %s", ColorYellow, ColorReset)
	branchName := gm.getUserInput()

	if branchName == "" {
		fmt.Printf("%s❌ Nom de branche requis!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	stashRef := fmt.Sprintf("stash@{%d}", index)
	output, err := gm.runGitCommand("stash", "branch", branchName, stashRef)
	if err != nil {
		fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
	} else {
		fmt.Printf("%s✅ Branche '%s' créée depuis le stash %s!%s\n", ColorGreen, branchName, stashRef, ColorReset)
	}
	gm.pause()
}

// Statistics
func (gm *GitManager) showGeneralStats() {
	fmt.Printf("%s%s📊 STATISTIQUES GÉNÉRALES%s\n", ColorBold, ColorBlue, ColorReset)
	fmt.Println(strings.Repeat("═", 40))

	totalCommits, _ := gm.runGitCommand("rev-list", "--count", "HEAD")
	fmt.Printf("%s📦 Total commits: %s%s%s\n", ColorGreen, ColorWhite, totalCommits, ColorReset)

	branches, _ := gm.runGitCommand("branch", "-r")
	branchCount := 0
	if branches != "" {
		branchCount = len(strings.Split(strings.TrimSpace(branches), "\n"))
	}
	localBranches, _ := gm.runGitCommand("branch")
	localCount := 0
	if localBranches != "" {
		localCount = len(strings.Split(strings.TrimSpace(localBranches), "\n"))
	}

	fmt.Printf("%s🌿 Branches locales: %s%d%s\n", ColorGreen, ColorWhite, localCount, ColorReset)
	fmt.Printf("%s🌐 Branches remote: %s%d%s\n", ColorGreen, ColorWhite, branchCount, ColorReset)

	tags, _ := gm.runGitCommand("tag")
	tagCount := 0
	if tags != "" {
		tagCount = len(strings.Split(strings.TrimSpace(tags), "\n"))
	}
	fmt.Printf("%s🏷️  Tags: %s%d%s\n", ColorGreen, ColorWhite, tagCount, ColorReset)

	firstCommit, _ := gm.runGitCommand("log", "--reverse", "--pretty=format:%h - %s (%an, %ad)", "--date=short", "-1")
	lastCommit, _ := gm.runGitCommand("log", "--pretty=format:%h - %s (%an, %ad)", "--date=short", "-1")

	fmt.Printf("\n%s📅 Premier commit:%s\n%s\n", ColorYellow, ColorReset, firstCommit)
	fmt.Printf("%s📅 Dernier commit:%s\n%s\n", ColorYellow, ColorReset, lastCommit)

	repoSize, _ := gm.runGitCommand("count-objects", "-vH")
	fmt.Printf("\n%s💾 Taille du dépôt:%s\n", ColorCyan, ColorReset)
	fmt.Println(repoSize)

	gm.pause()
}

func (gm *GitManager) showContributorStats() {
	fmt.Printf("%s%s👥 STATISTIQUES DES CONTRIBUTEURS%s\n", ColorBold, ColorBlue, ColorReset)
	fmt.Println(strings.Repeat("═", 45))

	contributors, _ := gm.runGitCommand("shortlog", "-sn", "--all")
	fmt.Printf("%s📈 Commits par contributeur:%s\n", ColorGreen, ColorReset)
	fmt.Println(contributors)

	fmt.Printf("\n%s📅 Activité par mois (12 derniers mois):%s\n", ColorYellow, ColorReset)
	activity, _ := gm.runGitCommand("log", "--pretty=format:%ad", "--date=format:%Y-%m", "--since=12.months.ago")
	if activity != "" {
		months := strings.Split(activity, "\n")
		monthCount := make(map[string]int)
		for _, month := range months {
			monthCount[month]++
		}

		var sortedMonths []string
		for month := range monthCount {
			sortedMonths = append(sortedMonths, month)
		}
		sort.Strings(sortedMonths)

		for _, month := range sortedMonths {
			fmt.Printf("%s: %d commits\n", month, monthCount[month])
		}
	}

	gm.pause()
}

func (gm *GitManager) showDetailedHistory() {
	fmt.Printf("%sNombre de commits à afficher (défaut: 20): %s", ColorYellow, ColorReset)
	countStr := gm.getUserInput()

	count := 20
	if countStr != "" {
		if c, err := strconv.Atoi(countStr); err == nil {
			count = c
		}
	}

	fmt.Println("\n1. Format court")
	fmt.Println("2. Format détaillé")
	fmt.Println("3. Format personnalisé")

	fmt.Printf("\n%sChoisissez un format: %s", ColorYellow, ColorReset)
	choice := gm.getUserInput()

	var output string
	var err error

	switch choice {
	case "1":
		output, err = gm.runGitCommand("log", "--oneline", "--graph", "--decorate", fmt.Sprintf("-%d", count))
	case "2":
		output, err = gm.runGitCommand("log", "--graph", "--pretty=format:%C(yellow)%h%C(reset) - %C(green)(%ar)%C(reset) %s %C(bold blue)<%an>%C(reset)", fmt.Sprintf("-%d", count))
	case "3":
		fmt.Printf("%sFormat personnalisé (ex: %%h - %%s (%%an)): %s", ColorYellow, ColorReset)
		format := gm.getUserInput()
		if format != "" {
			output, err = gm.runGitCommand("log", "--pretty=format:"+format, fmt.Sprintf("-%d", count))
		}
	default:
		output, err = gm.runGitCommand("log", "--oneline", fmt.Sprintf("-%d", count))
	}

	if err != nil {
		fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
	} else {
		fmt.Printf("%s📈 Historique:%s\n", ColorBlue, ColorReset)
		fmt.Println(output)
	}

	gm.pause()
}

func (gm *GitManager) showBranchGraph() {
	fmt.Printf("%s%s🌳 GRAPHIQUE DES BRANCHES%s\n", ColorBold, ColorBlue, ColorReset)
	fmt.Println(strings.Repeat("═", 35))

	fmt.Printf("%sNombre de commits à afficher (défaut: 30): %s", ColorYellow, ColorReset)
	countStr := gm.getUserInput()

	count := 30
	if countStr != "" {
		if c, err := strconv.Atoi(countStr); err == nil {
			count = c
		}
	}

	output, err := gm.runGitCommand("log", "--graph", "--pretty=format:%C(auto)%h%C(reset) - %C(green)%s%C(reset) %C(yellow)(%ar)%C(reset) %C(bold blue)<%an>%C(reset)", "--all", fmt.Sprintf("-%d", count))
	if err != nil {
		fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
	} else {
		fmt.Println(output)
	}

	gm.pause()
}

func (gm *GitManager) showFileStats() {
	for {
		gm.clearScreen()
		fmt.Printf("%s%s📁 STATISTIQUES PAR FICHIER%s\n", ColorBold, ColorBlue, ColorReset)
		fmt.Println(strings.Repeat("═", 35))

		fmt.Println("1. Fichiers les plus modifiés")
		fmt.Println("2. Lignes ajoutées/supprimées par fichier")
		fmt.Println("3. Historique d'un fichier spécifique")
		fmt.Println("0. Retour")

		fmt.Printf("\n%sChoisissez une option: %s", ColorYellow, ColorReset)
		choice := gm.getUserInput()

		switch choice {
		case "1":
			output, err := gm.runGitCommand("log", "--pretty=format:", "--name-only")
			if err != nil {
				fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
			} else {
				files := strings.Split(output, "\n")
				fileCount := make(map[string]int)
				for _, file := range files {
					if file != "" {
						fileCount[file]++
					}
				}

				type fileInfo struct {
					name  string
					count int
				}
				var sortedFiles []fileInfo
				for name, count := range fileCount {
					sortedFiles = append(sortedFiles, fileInfo{name, count})
				}

				sort.Slice(sortedFiles, func(i, j int) bool {
					return sortedFiles[i].count > sortedFiles[j].count
				})

				fmt.Printf("%s📊 Fichiers les plus modifiés:%s\n", ColorGreen, ColorReset)
				for i, file := range sortedFiles {
					if i >= 20 {
						break
					}
					fmt.Printf("%s%3d%s modifications - %s\n", ColorYellow, file.count, ColorReset, file.name)
				}
			}
			gm.pause()

		case "2":
			output, err := gm.runGitCommand("log", "--numstat", "--pretty=format:")
			if err != nil {
				fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
			} else {
				lines := strings.Split(output, "\n")
				fileStats := make(map[string][2]int) // [additions, deletions]

				for _, line := range lines {
					if line != "" {
						parts := strings.Fields(line)
						if len(parts) == 3 {
							if parts[0] == "-" { // Binary file
								continue
							}
							filename := parts[2]
							if additions, err := strconv.Atoi(parts[0]); err == nil {
								if deletions, err := strconv.Atoi(parts[1]); err == nil {
									stats := fileStats[filename]
									stats[0] += additions
									stats[1] += deletions
									fileStats[filename] = stats
								}
							}
						}
					}
				}

				fmt.Printf("%s📊 Lignes ajoutées/supprimées:%s\n", ColorGreen, ColorReset)
				for filename, stats := range fileStats {
					fmt.Printf("%s+%d%s %s-%d%s %s\n", ColorGreen, stats[0], ColorReset, ColorRed, stats[1], ColorReset, filename)
				}
			}
			gm.pause()

		case "3":
			fmt.Printf("%sNom du fichier: %s", ColorYellow, ColorReset)
			filename := gm.getUserInput()
			if filename != "" {
				output, err := gm.runGitCommand("log", "--follow", "--oneline", "--", filename)
				if err != nil {
					fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
				} else {
					fmt.Printf("%s📊 Historique de '%s':%s\n", ColorGreen, filename, ColorReset)
					fmt.Println(output)
				}
			}
			gm.pause()
		case "0":
			return
		default:
			fmt.Printf("%s❌ Option invalide!%s\n", ColorRed, ColorReset)
			gm.pause()
		}
	}
}

func (gm *GitManager) showBlame() {
	fmt.Printf("%sNom du fichier à analyser: %s", ColorYellow, ColorReset)
	filename := gm.getUserInput()

	if filename == "" {
		return
	}

	if _, err := os.Stat(filepath.Join(gm.currentPath, filename)); os.IsNotExist(err) {
		fmt.Printf("%s❌ Fichier '%s' introuvable!%s\n", ColorRed, filename, ColorReset)
		gm.pause()
		return
	}

	for {
		gm.clearScreen()
		fmt.Println("1. Blâme complet")
		fmt.Println("2. Blâme avec statistiques")
		fmt.Println("3. Blâme d'une plage de lignes")
		fmt.Println("0. Retour")

		fmt.Printf("\n%sChoisissez une option: %s", ColorYellow, ColorReset)
		choice := gm.getUserInput()

		var output string
		var err error

		switch choice {
		case "1":
			output, err = gm.runGitCommand("blame", filename)
		case "2":
			output, err = gm.runGitCommand("blame", "-w", "-C", "-C", "-C", filename)
		case "3":
			fmt.Printf("%sLigne de début: %s", ColorYellow, ColorReset)
			startLine := gm.getUserInput()
			fmt.Printf("%sLigne de fin: %s", ColorYellow, ColorReset)
			endLine := gm.getUserInput()

			if startLine != "" && endLine != "" {
				lineRange := fmt.Sprintf("-L%s,%s", startLine, endLine)
				output, err = gm.runGitCommand("blame", lineRange, filename)
			}
		case "0":
			return
		default:
			fmt.Printf("%s❌ Option invalide!%s\n", ColorRed, ColorReset)
			gm.pause()
			continue
		}

		if err != nil {
			fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
		} else {
			fmt.Printf("%s🔍 Blâme de '%s':%s\n", ColorBlue, filename, ColorReset)
			fmt.Println(output)
		}
		gm.pause()
	}
}

func (gm *GitManager) searchHistory() {
	for {
		gm.clearScreen()
		fmt.Println("1. Rechercher dans les messages de commit")
		fmt.Println("2. Rechercher dans le code (diff)")
		fmt.Println("3. Rechercher par auteur")
		fmt.Println("4. Rechercher par date")
		fmt.Println("0. Retour")

		fmt.Printf("\n%sChoisissez une option: %s", ColorYellow, ColorReset)
		choice := gm.getUserInput()

		var output string
		var err error
		var args []string

		switch choice {
		case "1":
			fmt.Printf("%sTerme à rechercher: %s", ColorYellow, ColorReset)
			term := gm.getUserInput()
			if term != "" {
				args = []string{"log", "--grep=" + term, "--oneline"}
				output, err = gm.runGitCommand(args...)
				fmt.Printf("%s🔍 Commits contenant '%s':%s\n", ColorBlue, term, ColorReset)
			}
		case "2":
			fmt.Printf("%sCode à rechercher: %s", ColorYellow, ColorReset)
			code := gm.getUserInput()
			if code != "" {
				args = []string{"log", "-S" + code, "--oneline"}
				output, err = gm.runGitCommand(args...)
				fmt.Printf("%s🔍 Commits modifiant '%s':%s\n", ColorBlue, code, ColorReset)
			}
		case "3":
			fmt.Printf("%sAuteur à rechercher: %s", ColorYellow, ColorReset)
			author := gm.getUserInput()
			if author != "" {
				args = []string{"log", "--author=" + author, "--oneline"}
				output, err = gm.runGitCommand(args...)
				fmt.Printf("%s🔍 Commits de '%s':%s\n", ColorBlue, author, ColorReset)
			}
		case "4":
			fmt.Printf("%sDate de début (YYYY-MM-DD): %s", ColorYellow, ColorReset)
			since := gm.getUserInput()
			fmt.Printf("%sDate de fin (YYYY-MM-DD, optionnel): %s", ColorYellow, ColorReset)
			until := gm.getUserInput()

			if since != "" {
				args = []string{"log", "--oneline", "--since=" + since}
				if until != "" {
					args = append(args, "--until="+until)
				}
				output, err = gm.runGitCommand(args...)
				fmt.Printf("%s🔍 Commits depuis %s", ColorBlue, since)
				if until != "" {
					fmt.Printf(" jusqu'à %s", until)
				}
				fmt.Printf(":%s\n", ColorReset)
			}
		case "0":
			return
		default:
			fmt.Printf("%s❌ Option invalide!%s\n", ColorRed, ColorReset)
			gm.pause()
			continue
		}

		if err != nil {
			fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
		} else if output == "" {
			fmt.Printf("%s❌ Aucun résultat%s\n", ColorRed, ColorReset)
		} else {
			fmt.Println(output)
		}
		gm.pause()
	}
}

// Tools
func (gm *GitManager) manageConfig() {
	for {
		gm.clearScreen()
		fmt.Println("1. Voir la configuration")
		fmt.Println("2. Configurer utilisateur")
		fmt.Println("3. Configurer email")
		fmt.Println("4. Autres configurations")
		fmt.Println("0. Retour")

		fmt.Printf("\n%sChoisissez une option: %s", ColorYellow, ColorReset)
		choice := gm.getUserInput()

		switch choice {
		case "1":
			output, err := gm.runGitCommand("config", "--list")
			if err != nil {
				fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
			} else {
				fmt.Printf("%s⚙️  Configuration actuelle:%s\n", ColorBlue, ColorReset)
				fmt.Println(output)
			}
			gm.pause()
		case "2":
			fmt.Printf("%sNom d'utilisateur: %s", ColorYellow, ColorReset)
			username := gm.getUserInput()
			if username != "" {
				output, err := gm.runGitCommand("config", "user.name", username)
				if err != nil {
					fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
				} else {
					fmt.Printf("%s✅ Nom d'utilisateur configuré!%s\n", ColorGreen, ColorReset)
				}
			}
			gm.pause()
		case "3":
			fmt.Printf("%sEmail: %s", ColorYellow, ColorReset)
			email := gm.getUserInput()
			if email != "" {
				output, err := gm.runGitCommand("config", "user.email", email)
				if err != nil {
					fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
				} else {
					fmt.Printf("%s✅ Email configuré!%s\n", ColorGreen, ColorReset)
				}
			}
			gm.pause()
		case "4":
			fmt.Println("Configurations rapides:")
			fmt.Println("a. Activer la couleur")
			fmt.Println("b. Configurer l'éditeur par défaut")
			fmt.Println("c. Configurer le push par défaut")

			fmt.Printf("\n%sChoisissez: %s", ColorYellow, ColorReset)
			subChoice := gm.getUserInput()

			switch subChoice {
			case "a":
				gm.runGitCommand("config", "color.ui", "auto")
				fmt.Printf("%s✅ Couleur activée!%s\n", ColorGreen, ColorReset)
			case "b":
				fmt.Printf("%sÉditeur (nano, vim, code, etc.): %s", ColorYellow, ColorReset)
				editor := gm.getUserInput()
				if editor != "" {
					gm.runGitCommand("config", "core.editor", editor)
					fmt.Printf("%s✅ Éditeur configuré!%s\n", ColorGreen, ColorReset)
				}
			case "c":
				gm.runGitCommand("config", "push.default", "simple")
				fmt.Printf("%s✅ Push par défaut configuré!%s\n", ColorGreen, ColorReset)
			}
			gm.pause()
		case "0":
			return
		default:
			fmt.Printf("%s❌ Option invalide!%s\n", ColorRed, ColorReset)
			gm.pause()
		}
	}
}

func (gm *GitManager) cleanRepo() {
	for {
		gm.clearScreen()
		fmt.Printf("%s🧹 NETTOYAGE DU DÉPÔT%s\n", ColorYellow, ColorReset)
		fmt.Println(strings.Repeat("═", 25))

		fmt.Println("1. Nettoyer les fichiers non trackés (clean)")
		fmt.Println("2. Nettoyer les objets inaccessibles (prune)")
		fmt.Println("3. Optimiser le dépôt (gc)")
		fmt.Println("4. Nettoyage complet")
		fmt.Println("0. Retour")

		fmt.Printf("\n%sChoisissez une option: %s", ColorYellow, ColorReset)
		choice := gm.getUserInput()

		switch choice {
		case "1":
			preview, _ := gm.runGitCommand("clean", "-n")
			if preview != "" {
				fmt.Printf("%sFichiers qui seront supprimés:%s\n", ColorRed, ColorReset)
				fmt.Println(preview)

				fmt.Printf("\n%sContinuer? (y/N): %s", ColorYellow, ColorReset)
				confirm := gm.getUserInput()
				if strings.ToLower(confirm) == "y" {
					output, err := gm.runGitCommand("clean", "-f")
					if err != nil {
						fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
					} else {
						fmt.Printf("%s✅ Fichiers non trackés supprimés!%s\n", ColorGreen, ColorReset)
					}
				}
			} else {
				fmt.Printf("%s✅ Aucun fichier à nettoyer!%s\n", ColorGreen, ColorReset)
			}
			gm.pause()
		case "2":
			fmt.Printf("%sNettoyage des objets inaccessibles...%s\n", ColorYellow, ColorReset)
			output, err := gm.runGitCommand("gc", "--prune=now")
			if err != nil {
				fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
			} else {
				fmt.Printf("%s✅ Objets inaccessibles nettoyés!%s\n", ColorGreen, ColorReset)
			}
			gm.pause()
		case "3":
			fmt.Printf("%sOptimisation du dépôt...%s\n", ColorYellow, ColorReset)
			output, err := gm.runGitCommand("gc", "--aggressive", "--prune=now")
			if err != nil {
				fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
			} else {
				fmt.Printf("%s✅ Dépôt optimisé!%s\n", ColorGreen, ColorReset)
			}
			gm.pause()
		case "4":
			fmt.Printf("%s⚠️  Nettoyage complet (peut prendre du temps). Continuer? (y/N): %s", ColorRed, ColorReset)
			confirm := gm.getUserInput()
			if strings.ToLower(confirm) == "y" {
				fmt.Printf("%sNettoyage en cours...%s\n", ColorYellow, ColorReset)
				gm.runGitCommand("clean", "-f", "-d")
				gm.runGitCommand("gc", "--aggressive", "--prune=now")
				gm.runGitCommand("reflog", "expire", "--expire=now", "--all")
				fmt.Printf("%s✅ Nettoyage complet terminé!%s\n", ColorGreen, ColorReset)
			}
			gm.pause()
		case "0":
			return
		default:
			fmt.Printf("%s❌ Option invalide!%s\n", ColorRed, ColorReset)
			gm.pause()
		}
	}
}

func (gm *GitManager) checkRepo() {
	for {
		gm.clearScreen()
		fmt.Printf("%s🔍 VÉRIFICATION DU DÉPÔT%s\n", ColorBlue, ColorReset)
		fmt.Println(strings.Repeat("═", 30))

		fmt.Println("1. Vérifier l'intégrité (fsck)")
		fmt.Println("2. Statistiques des objets (count-objects)")
		fmt.Println("0. Retour")

		fmt.Printf("\n%sChoisissez une option: %s", ColorYellow, ColorReset)
		choice := gm.getUserInput()

		switch choice {
		case "1":
			fmt.Printf("%sVérification de l'intégrité...%s\n", ColorYellow, ColorReset)
			output, err := gm.runGitCommand("fsck", "--full")
			if err != nil || output != "" {
				fmt.Printf("%s❌ Problèmes détectés:%s\n", ColorRed, ColorReset)
				fmt.Println(output)
			} else {
				fmt.Printf("%s✅ Intégrité vérifiée! Aucun problème trouvé.%s\n", ColorGreen, ColorReset)
			}
			gm.pause()
		case "2":
			fmt.Printf("%s📊 Statistiques détaillées:%s\n", ColorBlue, ColorReset)
			output, _ := gm.runGitCommand("count-objects", "-vH")
			fmt.Println(output)
			gm.pause()
		case "0":
			return
		default:
			fmt.Printf("%s❌ Option invalide!%s\n", ColorRed, ColorReset)
			gm.pause()
		}
	}
}

func (gm *GitManager) manageHooks() {
	hooksDir := filepath.Join(gm.currentPath, ".git", "hooks")

	for {
		gm.clearScreen()
		fmt.Printf("%s🪝 GESTION DES HOOKS%s\n", ColorPurple, ColorReset)
		fmt.Println(strings.Repeat("═", 25))

		hooks, err := os.ReadDir(hooksDir)
		if err != nil {
			fmt.Printf("%s❌ Impossible de lire le répertoire hooks%s\n", ColorRed, ColorReset)
			gm.pause()
			return
		}

		fmt.Printf("%s📋 Hooks disponibles:%s\n", ColorBlue, ColorReset)
		for _, hook := range hooks {
			if !strings.HasSuffix(hook.Name(), ".sample") {
				fmt.Printf("%s✅ %s%s (actif)\n", ColorGreen, hook.Name(), ColorReset)
			} else {
				fmt.Printf("%s⚪ %s%s (exemple)\n", ColorYellow, hook.Name(), ColorReset)
			}
		}

		fmt.Println("\n1. Activer un hook d'exemple")
		fmt.Println("2. Désactiver un hook")
		fmt.Println("3. Voir le contenu d'un hook")
		fmt.Println("4. Créer un hook personnalisé")
		fmt.Println("0. Retour")

		fmt.Printf("\n%sChoisissez une option: %s", ColorYellow, ColorReset)
		choice := gm.getUserInput()

		switch choice {
		case "1":
			fmt.Printf("%sNom du hook à activer (sans .sample): %s", ColorYellow, ColorReset)
			hookName := gm.getUserInput()
			if hookName != "" {
				samplePath := filepath.Join(hooksDir, hookName+".sample")
				hookPath := filepath.Join(hooksDir, hookName)

				if _, err := os.Stat(samplePath); err == nil {
					content, _ := os.ReadFile(samplePath)
					err := os.WriteFile(hookPath, content, 0755)
					if err != nil {
						fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, err.Error(), ColorReset)
					} else {
						fmt.Printf("%s✅ Hook '%s' activé!%s\n", ColorGreen, hookName, ColorReset)
					}
				} else {
					fmt.Printf("%s❌ Hook d'exemple '%s.sample' introuvable!%s\n", ColorRed, hookName, ColorReset)
				}
			}
			gm.pause()
		case "2":
			fmt.Printf("%sNom du hook à désactiver: %s", ColorYellow, ColorReset)
			hookName := gm.getUserInput()
			if hookName != "" {
				hookPath := filepath.Join(hooksDir, hookName)
				err := os.Remove(hookPath)
				if err != nil {
					fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, err.Error(), ColorReset)
				} else {
					fmt.Printf("%s✅ Hook '%s' désactivé!%s\n", ColorGreen, hookName, ColorReset)
				}
			}
			gm.pause()
		case "3":
			fmt.Printf("%sNom du hook à voir: %s", ColorYellow, ColorReset)
			hookName := gm.getUserInput()
			if hookName != "" {
				hookPath := filepath.Join(hooksDir, hookName)
				content, err := os.ReadFile(hookPath)
				if err != nil {
					fmt.Printf("%s❌ Hook introuvable: %s%s\n", ColorRed, err.Error(), ColorReset)
				} else {
					fmt.Printf("%s📄 Contenu du hook '%s':%s\n", ColorBlue, hookName, ColorReset)
					fmt.Println(string(content))
				}
			}
			gm.pause()
		case "4":
			fmt.Printf("%sNom du nouveau hook: %s", ColorYellow, ColorReset)
			hookName := gm.getUserInput()
			if hookName != "" {
				hookContent := `#!/bin/sh
# Hook personnalisé: ` + hookName + `
# Créé le ` + time.Now().Format("2006-01-02 15:04:05") + `

echo "Exécution du hook ` + hookName + `"

# Ajoutez votre code ici
exit 0
`
				hookPath := filepath.Join(hooksDir, hookName)
				err := os.WriteFile(hookPath, []byte(hookContent), 0755)
				if err != nil {
					fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, err.Error(), ColorReset)
				} else {
					fmt.Printf("%s✅ Hook '%s' créé!%s\n", ColorGreen, hookName, ColorReset)
				}
			}
			gm.pause()
		case "0":
			return
		default:
			fmt.Printf("%s❌ Option invalide!%s\n", ColorRed, ColorReset)
			gm.pause()
		}
	}
}

func (gm *GitManager) manageAliases() {
	// This is a complex feature, so we'll provide a simpler interface
	fmt.Printf("%sLes alias Git sont gérés dans votre fichier de configuration Git (~/.gitconfig).%s\n", ColorYellow, ColorReset)
	fmt.Println("Exemples d'alias utiles:")
	fmt.Println("  git config --global alias.co checkout")
	fmt.Println("  git config --global alias.br branch")
	fmt.Println("  git config --global alias.ci commit")
	fmt.Println("  git config --global alias.st status")
	fmt.Println("  git config --global alias.last 'log -1 HEAD'")
	gm.pause()
}

func (gm *GitManager) archiveRepo() {
	fmt.Printf("%sFormat d'archive (zip, tar.gz): %s", ColorYellow, ColorReset)
	format := gm.getUserInput()
	if format != "zip" && format != "tar.gz" {
		fmt.Printf("%s❌ Format invalide. Utilisez 'zip' ou 'tar.gz'.%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	defaultFileName := filepath.Base(gm.currentPath) + "." + format
	fmt.Printf("%sNom du fichier de sortie (défaut: %s): %s", ColorYellow, defaultFileName, ColorReset)
	outputFile := gm.getUserInput()
	if outputFile == "" {
		outputFile = defaultFileName
	}

	output, err := gm.runGitCommand("archive", "--format="+format, "-o", outputFile, "HEAD")
	if err != nil {
		fmt.Printf("%s❌ Erreur lors de la création de l'archive: %s%s\n", ColorRed, output, ColorReset)
	} else {
		fmt.Printf("%s✅ Archive '%s' créée avec succès!%s\n", ColorGreen, outputFile, ColorReset)
	}
	gm.pause()
}

// Main Menu Actions
func (gm *GitManager) changeDirectory() {
	fmt.Printf("%sChemin actuel: %s%s\n", ColorYellow, gm.currentPath, ColorReset)
	fmt.Printf("%sNouveau chemin du répertoire: %s", ColorYellow, ColorReset)
	newPath := gm.getUserInput()

	if newPath == "" {
		return
	}

	if strings.HasPrefix(newPath, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			newPath = filepath.Join(home, newPath[1:])
		}
	}

	if err := os.Chdir(newPath); err != nil {
		fmt.Printf("%s❌ Erreur lors du changement de répertoire: %s%s\n", ColorRed, err, ColorReset)
	} else {
		newWd, _ := os.Getwd()
		gm.currentPath = newWd
		fmt.Printf("%s✅ Répertoire changé pour: %s%s\n", ColorGreen, gm.currentPath, ColorReset)
	}
	gm.pause()
}

func (gm *GitManager) initRepo() {
	if gm.isGitRepo() {
		fmt.Printf("%s⚠️  Un dépôt Git existe déjà dans ce répertoire.%s\n", ColorYellow, ColorReset)
		gm.pause()
		return
	}

	fmt.Printf("%sInitialiser un nouveau dépôt Git ici? (y/N): %s", ColorYellow, ColorReset)
	confirm := gm.getUserInput()

	if strings.ToLower(confirm) == "y" {
		output, err := gm.runGitCommand("init")
		if err != nil {
			fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
		} else {
			fmt.Printf("%s✅ Dépôt Git initialisé!%s\n", ColorGreen, ColorReset)
			fmt.Println(output)
		}
	}
	gm.pause()
}

// NOUVELLES FONCTIONS SIMPLIFIÉES POUR L'ACCÈS RAPIDE

func (gm *GitManager) handleQuickCommit() {
	if !gm.isGitRepo() {
		fmt.Printf("%s❌ Ce répertoire n'est pas un dépôt Git!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	fmt.Printf("%s%s📦 COMMIT RAPIDE%s\n", ColorBold, ColorGreen, ColorReset)
	fmt.Println(strings.Repeat("═", 20))

	// Vérifier s'il y a des fichiers en stage
	staged, _ := gm.runGitCommand("diff", "--cached", "--name-only")
	if staged != "" {
		fmt.Printf("%s✅ Fichiers en stage:%s\n", ColorGreen, ColorReset)
		fmt.Println(staged)
		fmt.Printf("\n%sMessage de commit: %s", ColorYellow, ColorReset)
		message := gm.getUserInput()
		if message != "" {
			output, err := gm.runGitCommand("commit", "-m", message)
			if err != nil {
				fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
			} else {
				fmt.Printf("%s✅ Commit créé!%s\n", ColorGreen, ColorReset)
			}
		} else {
			fmt.Printf("%s❌ Message de commit vide. Annulation.%s\n", ColorRed, ColorReset)
		}
	} else {
		fmt.Printf("%s📋 Aucun fichier en stage.%s\n", ColorYellow, ColorReset)
		fmt.Printf("%s1.%s Voir l'historique des commits\n", ColorCyan, ColorReset)
		fmt.Printf("%s2.%s Aller au menu complet des commits\n", ColorCyan, ColorReset)
		fmt.Printf("%s0.%s Retour\n", ColorRed, ColorReset)

		fmt.Printf("\n%sChoisissez: %s", ColorYellow, ColorReset)
		choice := gm.getUserInput()

		switch choice {
		case "1":
			output, _ := gm.runGitCommand("log", "--oneline", "--graph", "--decorate", "-10")
			fmt.Printf("%s📈 10 derniers commits:%s\n", ColorBlue, ColorReset)
			fmt.Println(output)
			gm.pause()
		case "2":
			gm.handleCommitManagement()
			return // Ne pas faire de pause, on va dans le menu complet
		}
	}
	gm.pause()
}

func (gm *GitManager) handleQuickFiles() {
	if !gm.isGitRepo() {
		fmt.Printf("%s❌ Ce répertoire n'est pas un dépôt Git!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	fmt.Printf("%s%s📁 FICHIERS RAPIDE%s\n", ColorBold, ColorGreen, ColorReset)
	fmt.Println(strings.Repeat("═", 20))

	status := gm.getGitStatus()
	if status == "" {
		fmt.Printf("%s✅ Aucun changement détecté%s\n", ColorGreen, ColorReset)
		gm.pause()
		return
	}

	fmt.Printf("%s📋 Fichiers modifiés:%s\n", ColorYellow, ColorReset)
	lines := strings.Split(status, "\n")
	for _, line := range lines {
		if line != "" {
			gm.printFileStatus(line)
		}
	}

	fmt.Printf("\n%s1.%s Ajouter tous les fichiers et voir diff\n", ColorCyan, ColorReset)
	fmt.Printf("%s2.%s Voir seulement les différences\n", ColorCyan, ColorReset)
	fmt.Printf("%s3.%s Menu complet des fichiers\n", ColorCyan, ColorReset)
	fmt.Printf("%s0.%s Retour\n", ColorRed, ColorReset)

	fmt.Printf("\n%sChoisissez: %s", ColorYellow, ColorReset)
	choice := gm.getUserInput()

	switch choice {
	case "1":
		output, err := gm.runGitCommand("add", ".")
		if err != nil {
			fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
		} else {
			fmt.Printf("%s✅ Tous les fichiers ajoutés!%s\n", ColorGreen, ColorReset)
			// Afficher un résumé
			staged, _ := gm.runGitCommand("diff", "--cached", "--stat")
			fmt.Printf("%s📊 Résumé des changements:%s\n", ColorBlue, ColorReset)
			fmt.Println(staged)
		}
		gm.pause()
	case "2":
		output, _ := gm.runGitCommand("diff", "--color=always")
		if output != "" {
			fmt.Printf("%s📊 Différences:%s\n", ColorBlue, ColorReset)
			fmt.Println(output)
		} else {
			fmt.Printf("%s✅ Aucun changement non stagé à afficher.%s\n", ColorGreen, ColorReset)
		}
		gm.pause()
	case "3":
		gm.handleFileManagement()
		return
	}
}

func (gm *GitManager) handleQuickBranch() {
	if !gm.isGitRepo() {
		fmt.Printf("%s❌ Ce répertoire n'est pas un dépôt Git!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	currentBranch := gm.getCurrentBranch()
	fmt.Printf("%s%s🌿 BRANCHES RAPIDE%s\n", ColorBold, ColorGreen, ColorReset)
	fmt.Println(strings.Repeat("═", 20))
	fmt.Printf("%sBranche actuelle: %s%s%s\n\n", ColorBlue, ColorCyan, currentBranch, ColorReset)

	branches, _ := gm.runGitCommand("branch", "--format=%(refname:short)")
	fmt.Printf("%sBranches locales:%s\n", ColorYellow, ColorReset)
	fmt.Println(branches)

	fmt.Printf("\n%s1.%s Créer une nouvelle branche\n", ColorCyan, ColorReset)
	fmt.Printf("%s2.%s Changer de branche\n", ColorCyan, ColorReset)
	fmt.Printf("%s3.%s Menu complet des branches\n", ColorCyan, ColorReset)
	fmt.Printf("%s0.%s Retour\n", ColorRed, ColorReset)

	fmt.Printf("\n%sChoisissez: %s", ColorYellow, ColorReset)
	choice := gm.getUserInput()

	switch choice {
	case "1":
		gm.createBranch()
	case "2":
		gm.switchBranch()
	case "3":
		gm.handleBranchManagement()
		return
	}
}

func (gm *GitManager) handleQuickRemote() {
	if !gm.isGitRepo() {
		fmt.Printf("%s❌ Ce répertoire n'est pas un dépôt Git!%s\n", ColorRed, ColorReset)
		gm.pause()
		return
	}

	fmt.Printf("%s%s🔄 REMOTE RAPIDE%s\n", ColorBold, ColorGreen, ColorReset)
	fmt.Println(strings.Repeat("═", 20))

	// Vérifier s'il y a des commits en avance/retard
	// Fetch pour s'assurer que les informations sont à jour
	gm.runGitCommand("fetch", "origin") // Fetch silently to update remote tracking branches
	ahead, _ := gm.runGitCommand("rev-list", "--count", "@{u}..HEAD")
	behind, _ := gm.runGitCommand("rev-list", "--count", "HEAD..@{u}")

	if ahead != "0" && ahead != "" {
		fmt.Printf("%s📤 %s commit(s) à pusher%s\n", ColorYellow, ahead, ColorReset)
	}
	if behind != "0" && behind != "" {
		fmt.Printf("%s📥 %s commit(s) à puller%s\n", ColorYellow, behind, ColorReset)
	}
	if (ahead == "0" || ahead == "") && (behind == "0" || behind == "") {
		fmt.Printf("%s✅ Votre branche est à jour avec le remote.%s\n", ColorGreen, ColorReset)
	}

	fmt.Printf("\n%s1.%s Push rapide (origin + branche actuelle)\n", ColorCyan, ColorReset)
	fmt.Printf("%s2.%s Pull rapide (origin + branche actuelle)\n", ColorCyan, ColorReset)
	fmt.Printf("%s3.%s Status remote complet (fetch)\n", ColorCyan, ColorReset) // Renommé pour plus de clarté
	fmt.Printf("%s4.%s Menu complet des remotes\n", ColorCyan, ColorReset)
	fmt.Printf("%s0.%s Retour\n", ColorRed, ColorReset)

	fmt.Printf("\n%sChoisissez: %s", ColorYellow, ColorReset)
	choice := gm.getUserInput()

	currentBranch := gm.getCurrentBranch()

	switch choice {
	case "1":
		fmt.Printf("%sPush vers origin/%s...%s\n", ColorYellow, currentBranch, ColorReset)
		output, err := gm.runGitCommand("push", "origin", currentBranch)
		if err != nil {
			fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
		} else {
			fmt.Printf("%s✅ Push terminé!%s\n", ColorGreen, ColorReset)
		}
		gm.pause()
	case "2":
		fmt.Printf("%sPull depuis origin/%s...%s\n", ColorYellow, currentBranch, ColorReset)
		output, err := gm.runGitCommand("pull", "origin", currentBranch)
		if err != nil {
			fmt.Printf("%s❌ Erreur: %s%s\n", ColorRed, output, ColorReset)
		} else {
			fmt.Printf("%s✅ Pull terminé!%s\n", ColorGreen, ColorReset)
		}
		gm.pause()
	case "3":
		gm.fetchFromRemote() // This function already pauses
	case "4":
		gm.handleRemoteManagement()
		return
	}
}

// Modification de la fonction main pour gérer les raccourcis
func main() {
	gm := NewGitManager()
	for {
		gm.clearScreen()
		gm.showMenu()
		choice := strings.ToUpper(gm.getUserInput()) // Convertir en majuscule pour les raccourcis

		switch choice {
		// RACCOURCIS ACCÈS RAPIDE
		case "S":
			gm.handleDetailedStatus()
		case "C":
			gm.handleQuickCommit() // Nouvelle fonction simplifiée
		case "F":
			gm.handleQuickFiles() // Nouvelle fonction simplifiée
		case "B":
			gm.handleQuickBranch() // Nouvelle fonction simplifiée
		case "R":
			gm.handleQuickRemote() // Nouvelle fonction simplifiée

		// MENU COMPLET (existant)
		case "1":
			gm.handleDetailedStatus()
		case "2":
			gm.handleBranchManagement()
		case "3":
			gm.handleCommitManagement()
		case "4":
			gm.handleRemoteManagement()
		case "5":
			gm.handleFileManagement()
		case "6":
			gm.handleTagManagement()
		case "7":
			gm.handleStashManagement()
		case "8":
			gm.handleStatistics()
		case "9":
			gm.handleTools()
		case "10":
			gm.changeDirectory()
		case "11":
			gm.initRepo()
		case "0":
			fmt.Println("👋 Au revoir!")
			return
		default:
			fmt.Printf("%s❌ Option invalide! Utilisez les chiffres (0-11) ou les lettres (S,C,F,B,R)%s\n", ColorRed, ColorReset)
			gm.pause()
		}
	}
}