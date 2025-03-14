package pipeline

import (
	"fmt"
	"regexp"
	"strings"
	"yaylog/internal/config"
	"yaylog/internal/consts"
	"yaylog/internal/pkgdata"
)

type (
	PackageInfo     = pkgdata.PackageInfo
	FilterCondition = pkgdata.FilterCondition
)

var packageListRegex = regexp.MustCompile(`^([a-z0-9][a-z0-9._-]*[a-z0-9])(,([a-z0-9][a-z0-9._-]*[a-z0-9]))*$`)

func PreprocessFiltering(
	cfg config.Config,
	packages []pkgdata.PackageInfo,
	reportProgress pkgdata.ProgressReporter,
) ([]pkgdata.PackageInfo, error) {
	if len(cfg.FilterQueries) == 0 {
		return packages, nil
	}

	filterConditions, err := queriesToConditions(cfg.FilterQueries)
	if err != nil {
		return []pkgdata.PackageInfo{}, err
	}

	return pkgdata.FilterPackages(packages, filterConditions, reportProgress), nil
}

func queriesToConditions(filterQueries map[consts.FieldType]string) ([]pkgdata.FilterCondition, error) {
	conditions := make([]pkgdata.FilterCondition, 0, len(filterQueries))

	for fieldType, value := range filterQueries {
		var condition pkgdata.FilterCondition
		var err error

		switch fieldType {
		case consts.FieldDate:
			condition, err = parseDateFilterCondition(value)
		case consts.FieldSize:
			condition, err = parseSizeFilterCondition(value)
		case consts.FieldName, consts.FieldRequiredBy, consts.FieldDepends, consts.FieldProvides:
			condition, err = parsePackageFilterCondition(fieldType, value)
		case consts.FieldReason:
			condition, err = parseReasonFilterCondition(value)
		default:
			err = fmt.Errorf("unsupported filter type: %s", fieldType)
		}

		if err != nil {
			return []pkgdata.FilterCondition{}, err
		}

		conditions = append(conditions, condition)
	}

	return conditions, nil
}

func parsePackageFilterCondition(
	fieldType consts.FieldType,
	packageListInput string,
) (FilterCondition, error) {
	if !packageListRegex.MatchString(packageListInput) {
		return FilterCondition{}, fmt.Errorf("invalid package list: %s", packageListInput)
	}

	packageNames := strings.Split(packageListInput, ",")
	return NewPackageCondition(fieldType, packageNames)
}

func parseReasonFilterCondition(installReason string) (FilterCondition, error) {
	if installReason != config.ReasonExplicit && installReason != config.ReasonDependency {
		return FilterCondition{}, fmt.Errorf("invalid install reason filter: %s", installReason)
	}

	return NewReasonCondition(installReason), nil
}

func parseDateFilterCondition(value string) (FilterCondition, error) {
	dateFilter, err := parseDateFilter(value)
	if err != nil {
		return pkgdata.FilterCondition{}, fmt.Errorf("invalid date filter: %v", err)
	}

	if err = validateDateFilter(dateFilter); err != nil {
		return pkgdata.FilterCondition{}, err
	}

	return NewDateCondition(dateFilter), nil
}

func parseSizeFilterCondition(value string) (FilterCondition, error) {
	sizeFilter, err := parseSizeFilter(value)
	if err != nil {
		return FilterCondition{}, fmt.Errorf("invalid size filter: %v", err)
	}

	if err = validateSizeFilter(sizeFilter); err != nil {
		return FilterCondition{}, err
	}

	return NewSizeCondition(sizeFilter), nil
}
