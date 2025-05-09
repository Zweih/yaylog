package phase

import (
	"fmt"
	"slices"
	"yaylog/internal/config"
	"yaylog/internal/consts"
	out "yaylog/internal/display"
	"yaylog/internal/pipeline/filtering"
	"yaylog/internal/pipeline/meta"
	"yaylog/internal/pkgdata"
)

func LoadCacheStep(
	_ config.Config,
	_ []*PkgInfo,
	_ ProgressReporter,
	pipelineCtx *meta.PipelineContext,
) ([]*PkgInfo, error) {
	pkgPtrs, err := pkgdata.LoadProtoCache()
	if err == nil {
		pipelineCtx.UsedCache = true
	}

	// TODO: use ProgressReporter to report cache status
	return pkgPtrs, nil
}

// TODO: add progress reporting
func FetchStep(
	_ config.Config,
	pkgPtrs []*PkgInfo,
	_ ProgressReporter,
	pipelineCtx *meta.PipelineContext,
) ([]*PkgInfo, error) {
	if !pipelineCtx.UsedCache {
		var err error
		pkgPtrs, err = pkgdata.FetchPackages()
		if err != nil {
			out.WriteLine(fmt.Sprintf(
				"Warning: Some packages may be missing due to corrupted package database: %v",
				err,
			))
		}
	}

	return pkgPtrs, nil
}

func ReverseDepStep(
	cfg config.Config,
	pkgPtrs []*pkgdata.PkgInfo,
	reportProgress ProgressReporter,
	pipelineCtx *meta.PipelineContext,
) ([]*PkgInfo, error) {
	if pipelineCtx.UsedCache {
		return pkgPtrs, nil
	}

	_, hasRequiredByFilter := cfg.FilterQueries[consts.FieldRequiredBy]
	hasRequiredByField := slices.Contains(cfg.Fields, consts.FieldRequiredBy)

	if !hasRequiredByField && !hasRequiredByFilter {
		return pkgPtrs, nil
	}

	return pkgdata.CalculateReverseDependencies(pkgPtrs, reportProgress)
}

// TODO: add progress reporting
func SaveCacheStep(
	_ config.Config,
	pkgPtrs []*PkgInfo,
	_ ProgressReporter,
	pipelineCtx *meta.PipelineContext,
) ([]*PkgInfo, error) {
	if !pipelineCtx.UsedCache {
		// TODO: we can probably save the file concurrently
		err := pkgdata.SaveProtoCache(pkgPtrs)
		if err != nil {
			out.WriteLine(fmt.Sprintf("Warning: Error saving cache: %v", err))
		}
	}

	return pkgPtrs, nil
}

func FilterStep(
	cfg config.Config,
	pkgPtrs []*PkgInfo,
	reportProgress ProgressReporter,
	_ *meta.PipelineContext,
) ([]*PkgInfo, error) {
	if len(cfg.FilterQueries) == 0 {
		return pkgPtrs, nil
	}

	filterConditions, err := filtering.QueriesToConditions(cfg.FilterQueries)
	if err != nil {
		return []*pkgdata.PkgInfo{}, err
	}

	return pkgdata.FilterPackages(pkgPtrs, filterConditions, reportProgress), nil
}

func SortStep(
	cfg config.Config,
	pkgPtrs []*PkgInfo,
	reportProgress ProgressReporter,
	_ *meta.PipelineContext,
) ([]*PkgInfo, error) {
	comparator := pkgdata.GetComparator(cfg.SortOption.Field, cfg.SortOption.Asc)
	phase := "Sorting packages"

	// threshold is 500 as that is where merge sorting chunk performance overtakes timsort
	if len(pkgPtrs) < pkgdata.ConcurrentSortThreshold {
		return pkgdata.SortNormally(pkgPtrs, comparator, phase, reportProgress), nil
	}

	return pkgdata.SortConcurrently(pkgPtrs, comparator, phase, reportProgress), nil
}
