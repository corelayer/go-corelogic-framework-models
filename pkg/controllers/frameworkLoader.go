package controllers

import (
	"fmt"
	"github.com/corelayer/go-corelogic-framework-models/pkg/models"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

type FrameworkLoader struct {
}

//func (c *FrameworkLoader) ListAvailableVersions() ([]string, error) {
//	var output []string
//
//	dirEntries, err := c.Assets.ReadDir("assets/framework")
//	if err != nil {
//		log.Fatal(err)
//		return output, err
//	}
//
//	for _, dirEntry := range dirEntries {
//		output = append(output, dirEntry.Name())
//	}
//	return output, nil
//}

//func (c *FrameworkLoader) Load(version string) (FrameworkController, error) {
//	var frameworkController FrameworkController
//	var err error
//	release, err := parseVersion(strings.Split(version, "."))
//	log.Printf("Release %v\n", release)
//
//	frameworkController.Release = release
//
//	releases, err := c.ListPreviousVersionsForMajorRelease(release)
//	if err != nil {
//		return frameworkController, err
//	}
//
//	frameworkController.Frameworks, err = c.LoadPreviousVersions(releases)
//	return frameworkController, err
//}

//func (c *FrameworkLoader) ListPreviousVersionsForMajorRelease(release models.Release) ([]models.Release, error) {
//	var output []models.Release
//
//	dirEntries, err := c.Assets.ReadDir("assets/framework")
//	if err != nil {
//		log.Fatal(err)
//		return output, err
//	}
//
//	for _, dirEntry := range dirEntries {
//		dirName := dirEntry.Name()
//		var currentRelease, loopErr = parseVersion(strings.Split(dirName, "."))
//		if loopErr != nil {
//			return output, loopErr
//		}
//
//		if currentRelease.Major == release.Major && strings.Compare(currentRelease.GetSemanticVersion(), release.GetSemanticVersion()) <= 0 {
//			log.Printf("Adding release %s to list of previous versions for %s", currentRelease.GetSemanticVersion(), release.GetSemanticVersion())
//			output = append(output, currentRelease)
//		}
//	}
//
//	return output, nil
//}

//func (c *FrameworkLoader) LoadPreviousVersions(releases []models.Release) (map[string]models.Framework, error) {
//	output := make(map[string]models.Framework)
//
//	for _, r := range releases {
//		currentFramework, loadErr := c.LoadVersion(r.GetSemanticVersion())
//		if loadErr != nil {
//			return output, loadErr
//		}
//		output[r.GetSemanticVersion()] = currentFramework
//	}
//
//	return output, nil
//}

//func parseVersion(version []string) (models.Release, error) {
//	major := 0
//	minor := 0
//	patch := 0
//	var err error
//
//	if len(version) != 3 {
//		err = fmt.Errorf("invalid input: %v", version)
//		return models.Release{
//			Major: 0,
//			Minor: 0,
//			Patch: 0,
//		}, err
//
//	}
//
//	major, err = strconv.Atoi(version[0])
//	if err != nil {
//		return models.Release{
//			Major: major,
//			Minor: minor,
//			Patch: patch,
//		}, err
//	}
//	minor, err = strconv.Atoi(version[1])
//	if err != nil {
//		return models.Release{
//			Major: major,
//			Minor: minor,
//			Patch: patch,
//		}, err
//	}
//	patch, err = strconv.Atoi(version[2])
//	if err != nil {
//		return models.Release{
//			Major: major,
//			Minor: minor,
//			Patch: patch,
//		}, err
//	}
//
//	return models.Release{
//		Major: major,
//		Minor: minor,
//		Patch: patch,
//	}, err
//
//}

func (c *FrameworkLoader) LoadFromDisk(rootDir string) (models.Framework, error) {
	//defer general.FinishTimer(general.StartTimer("Loading framework " + version))
	framework := models.Framework{}
	var source []byte
	var err error

	//log.Printf("Reading framework file at %s\n", rootDir+"/framework.yaml")
	source, err = os.ReadFile(rootDir + "/framework.yaml")
	if err != nil {
		fmt.Println(source)
		return framework, err
	}

	err = yaml.Unmarshal(source, &framework)
	if err != nil {
		log.Fatal(err)
		return framework, err
	}

	framework.Packages = []models.Package{}

	subDirs, err := os.ReadDir(rootDir + "/packages")
	if err != nil {
		log.Fatal(err)
		return framework, err
	}

	for _, d := range subDirs {
		if d.IsDir() {
			var p models.Package
			p, err = c.getPackagesFromDirectory(rootDir, d.Name())
			if err != nil {
				return framework, err
			}
			framework.Packages = append(framework.Packages, p)
		}
	}

	return framework, err
}

func (c *FrameworkLoader) getPackagesFromDirectory(rootDir string, directoryName string) (models.Package, error) {
	// defer general.FinishTimer(general.StartTimer("GetPackagesFromDirectory " + rootDir + "/packages/" + directoryName))

	myPackage := models.Package{
		Name:    directoryName,
		Modules: []models.Module{},
	}

	files, err := os.ReadDir(rootDir + "/packages/" + myPackage.Name)
	if err != nil {
		log.Fatal(err)
		return myPackage, err
	}

	for _, f := range files {
		if !f.IsDir() {
			if filepath.Ext(f.Name()) == ".yaml" {
				// log.Println(f.Name())
				var module models.Module
				module, err = c.getModuleFromFile(rootDir + "/packages/" + myPackage.Name + "/" + f.Name())
				if err != nil {
					return myPackage, err
				}
				myPackage.Modules = append(myPackage.Modules, module)
			}
		} else {
			var modules []models.Module
			modules, err = c.getModulesFromDirectory(rootDir + "/packages/" + myPackage.Name + "/" + f.Name())
			if err != nil {
				return myPackage, err
			}
			myPackage.Modules = append(myPackage.Modules, modules...)
		}
	}
	return myPackage, err
}

func (c *FrameworkLoader) getModuleFromFile(filePath string) (models.Module, error) {
	// defer general.FinishTimer(general.StartTimer("GetModuleFromFile " + filePath))

	module := models.Module{}

	moduleSource, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
		return module, err
	}

	err = yaml.Unmarshal(moduleSource, &module)
	if err != nil {
		log.Fatal(err)
	}

	return module, err
}

func (c *FrameworkLoader) getModulesFromDirectory(filePath string) ([]models.Module, error) {
	// defer general.FinishTimer(general.StartTimer("GetModulesFromDirectory " + filePath))

	var modules []models.Module

	files, err := os.ReadDir(filePath)
	if err != nil {
		log.Fatal(err)
		return modules, err
	}

	for _, f := range files {
		if !f.IsDir() {
			if filepath.Ext(f.Name()) == ".yaml" {
				// log.Println(f.Name())
				module, err := c.getModuleFromFile(filePath + "/" + f.Name())
				if err != nil {
					log.Fatal(err)
					return modules, err
				}
				modules = append(modules, module)
			}
		}
	}

	return modules, err
}
