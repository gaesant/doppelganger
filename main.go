package main

import (
	"doppelganger/functions"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

var (
	appData string // Windows AppData/Local folder path
	roaming string // Windows AppData/roaming folder
	appAsar string // Discord "app.asar" file location

	version string // Discord version

	basePath string // Discord base-path

	// Our patch's target
	patchTarget string = `return _path.default.join(userDataRoot, 'discord' + (buildInfo.releaseChannel == 'stable' ? '' : buildInfo.releaseChannel));`
)

func main() {

	fmt.Println("O patch que você está prestes a aplicar modificará a pasta-raiz do seu Discord podendo")
	fmt.Println("causar a completa inutilização do programa, portanto, verifique no repositório oficial: https://github.com/m3gadrive/doppelganger")
	fmt.Println("se esta é realmente a última versão deste patch, e se este patch ainda é funcional.")

	fmt.Println("Caso saiba plenamente o que está fazendo, pressione qualquer tecla. Se não, feche o programa.")

	fmt.Scanln()

	fmt.Println("Iniciando aplicação...")
	appData, _ = os.UserCacheDir()
	basePath = appData + "/Discord"

	version = getVersion()

	fmt.Println("Discord localizado na versão: ", strings.ReplaceAll(version, "app-", ""))

	appAsar = path.Join(basePath, version, "resources", "app.asar")

	archive, err := functions.Decode(appAsar)

	fmt.Println("Instância do app.asar decodificada, aplicando modificações...")

	if err != nil {
		handleErrors("Não encontrar o arquivo, por favor certifique-se de que o Discord está instalado.")
		fmt.Scanln()
		return
	}

	var hashedPath = functions.RandomString(18)
	functions.Modify(archive, patchTarget, strings.ReplaceAll(patchTarget, `'discord'`, hashedPath))

	roaming, err = os.UserConfigDir()
	if err != nil {
		handleErrors("Falha ao buscar o diretório AppData/Roaming. ")
		return
	}

	err = os.Rename(path.Join(roaming, "testtest"), hashedPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Modificações aplicadas, falta pouquinho! Vou salvar os arquivos.")

	err = functions.EncodeTo(archive, appAsar)

	if err != nil {
		handleErrors("Oops! Houve um erro enquanto tentava salvar um arquivo! Por favor, execute o programa como administrador e tente novamente. ")
		return
	}

	fmt.Println("Patch aplicado com sucesso ;)")
	fmt.Println("Pressione qualquer tecla para sair =D")

	fmt.Scanln()
}

func getVersion() string {
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		log.Fatal(err)
	}

	var appVersion string

	for _, f := range files {
		fileName := f.Name()

		if strings.Contains(fileName, "app-") {
			appVersion = fileName
		}
	}

	return appVersion
}

func handleErrors(message string) {
	fmt.Println(message)
	fmt.Println("Pressione qualquer tecla para sair.")
	fmt.Scanln()
	return
}
