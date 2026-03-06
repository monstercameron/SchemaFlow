package fluent

// AnnotateRequest is a fluent builder for Annotate.
type AnnotateRequest[T any] struct {
	commonRequest[AnnotateRequest[T], AnnotateOptions]
	input T
}

func newAnnotateRequest[T any](input T, opts AnnotateOptions) AnnotateRequest[T] {
	return AnnotateRequest[T]{
		commonRequest: commonRequest[AnnotateRequest[T], AnnotateOptions]{
			opts: opts,
			lift: func(next AnnotateOptions) AnnotateRequest[T] {
				return newAnnotateRequest(input, next)
			},
			mutate: func(current AnnotateOptions, fn func(CommonOptions) CommonOptions) AnnotateOptions {
				current.CommonOptions = fn(current.CommonOptions)
				return current
			},
		},
		input: input,
	}
}

// Annotating starts a fluent Annotate request.
func Annotating[T any](input T) AnnotateRequest[T] {
	return newAnnotateRequest(input, NewAnnotateOptions())
}

func (r AnnotateRequest[T]) Types(annotationTypes ...string) AnnotateRequest[T] {
	opts := r.opts
	opts.AnnotationTypes = append([]string(nil), annotationTypes...)
	return r.WithOptions(opts)
}

func (r AnnotateRequest[T]) Run() (AnnotateResult, error) {
	return Annotate[T](r.input, r.opts)
}

// ClusterRequest is a fluent builder for Cluster.
type ClusterRequest[T any] struct {
	commonRequest[ClusterRequest[T], ClusterOptions]
	items []T
}

func newClusterRequest[T any](items []T, opts ClusterOptions) ClusterRequest[T] {
	return ClusterRequest[T]{
		commonRequest: commonRequest[ClusterRequest[T], ClusterOptions]{
			opts: opts,
			lift: func(next ClusterOptions) ClusterRequest[T] {
				return newClusterRequest(items, next)
			},
			mutate: func(current ClusterOptions, fn func(CommonOptions) CommonOptions) ClusterOptions {
				current.CommonOptions = fn(current.CommonOptions)
				return current
			},
		},
		items: items,
	}
}

// Clustering starts a fluent Cluster request.
func Clustering[T any](items []T) ClusterRequest[T] {
	return newClusterRequest(items, NewClusterOptions())
}

func (r ClusterRequest[T]) By(criteria string) ClusterRequest[T] {
	opts := r.opts
	opts.ClusterBy = criteria
	return r.WithOptions(opts)
}

func (r ClusterRequest[T]) Clusters(n int) ClusterRequest[T] {
	opts := r.opts
	opts.NumClusters = n
	return r.WithOptions(opts)
}

func (r ClusterRequest[T]) Run() (ClusterResult[T], error) {
	return Cluster[T](r.items, r.opts)
}

// RankRequest is a fluent builder for Rank.
type RankRequest[T any] struct {
	commonRequest[RankRequest[T], RankOptions]
	items []T
}

func newRankRequest[T any](items []T, opts RankOptions) RankRequest[T] {
	return RankRequest[T]{
		commonRequest: commonRequest[RankRequest[T], RankOptions]{
			opts: opts,
			lift: func(next RankOptions) RankRequest[T] {
				return newRankRequest(items, next)
			},
			mutate: func(current RankOptions, fn func(CommonOptions) CommonOptions) RankOptions {
				current.CommonOptions = fn(current.CommonOptions)
				return current
			},
		},
		items: items,
	}
}

// Ranking starts a fluent Rank request.
func Ranking[T any](items []T) RankRequest[T] {
	return newRankRequest(items, NewRankOptions())
}

func (r RankRequest[T]) By(query string) RankRequest[T] {
	opts := r.opts
	opts.Query = query
	return r.WithOptions(opts)
}

func (r RankRequest[T]) Top(n int) RankRequest[T] {
	opts := r.opts
	opts.TopK = n
	return r.WithOptions(opts)
}

func (r RankRequest[T]) MinScore(score float64) RankRequest[T] {
	opts := r.opts
	opts.MinScore = score
	return r.WithOptions(opts)
}

func (r RankRequest[T]) Run() (RankResult[T], error) {
	return Rank[T](r.items, r.opts)
}

// CompressRequest is a fluent builder for Compress.
type CompressRequest[T any] struct {
	commonRequest[CompressRequest[T], CompressOptions]
	input T
}

func newCompressRequest[T any](input T, opts CompressOptions) CompressRequest[T] {
	return CompressRequest[T]{
		commonRequest: commonRequest[CompressRequest[T], CompressOptions]{
			opts: opts,
			lift: func(next CompressOptions) CompressRequest[T] {
				return newCompressRequest(input, next)
			},
			mutate: func(current CompressOptions, fn func(CommonOptions) CommonOptions) CompressOptions {
				current.CommonOptions = fn(current.CommonOptions)
				return current
			},
		},
		input: input,
	}
}

// Compressing starts a fluent Compress request.
func Compressing[T any](input T) CompressRequest[T] {
	return newCompressRequest(input, NewCompressOptions())
}

func (r CompressRequest[T]) Run() (CompressResult[T], error) {
	return Compress[T](r.input, r.opts)
}

// CompressTextRequest is a fluent builder for CompressText.
type CompressTextRequest struct {
	commonRequest[CompressTextRequest, CompressOptions]
	input string
}

func newCompressTextRequest(input string, opts CompressOptions) CompressTextRequest {
	return CompressTextRequest{
		commonRequest: commonRequest[CompressTextRequest, CompressOptions]{
			opts: opts,
			lift: func(next CompressOptions) CompressTextRequest {
				return newCompressTextRequest(input, next)
			},
			mutate: func(current CompressOptions, fn func(CommonOptions) CommonOptions) CompressOptions {
				current.CommonOptions = fn(current.CommonOptions)
				return current
			},
		},
		input: input,
	}
}

// CompressingText starts a fluent CompressText request.
func CompressingText(input string) CompressTextRequest {
	return newCompressTextRequest(input, NewCompressOptions())
}

func (r CompressTextRequest) Run() (string, error) {
	return CompressText(r.input, r.opts)
}

// DecomposeRequest is a fluent builder for Decompose.
type DecomposeRequest[T any] struct {
	commonRequest[DecomposeRequest[T], DecomposeOptions]
	input T
}

func newDecomposeRequest[T any](input T, opts DecomposeOptions) DecomposeRequest[T] {
	return DecomposeRequest[T]{
		commonRequest: commonRequest[DecomposeRequest[T], DecomposeOptions]{
			opts: opts,
			lift: func(next DecomposeOptions) DecomposeRequest[T] {
				return newDecomposeRequest(input, next)
			},
			mutate: func(current DecomposeOptions, fn func(CommonOptions) CommonOptions) DecomposeOptions {
				current.CommonOptions = fn(current.CommonOptions)
				return current
			},
		},
		input: input,
	}
}

// Decomposing starts a fluent Decompose request.
func Decomposing[T any](input T) DecomposeRequest[T] {
	return newDecomposeRequest(input, NewDecomposeOptions())
}

func (r DecomposeRequest[T]) Run() (DecomposeResult[T], error) {
	return Decompose[T](r.input, r.opts)
}

// DecomposeSliceRequest is a fluent builder for DecomposeToSlice.
type DecomposeSliceRequest[T any, U any] struct {
	commonRequest[DecomposeSliceRequest[T, U], DecomposeOptions]
	input T
}

func newDecomposeSliceRequest[T any, U any](input T, opts DecomposeOptions) DecomposeSliceRequest[T, U] {
	return DecomposeSliceRequest[T, U]{
		commonRequest: commonRequest[DecomposeSliceRequest[T, U], DecomposeOptions]{
			opts: opts,
			lift: func(next DecomposeOptions) DecomposeSliceRequest[T, U] {
				return newDecomposeSliceRequest[T, U](input, next)
			},
			mutate: func(current DecomposeOptions, fn func(CommonOptions) CommonOptions) DecomposeOptions {
				current.CommonOptions = fn(current.CommonOptions)
				return current
			},
		},
		input: input,
	}
}

// DecomposingInto starts a fluent DecomposeToSlice request.
func DecomposingInto[T any, U any](input T) DecomposeSliceRequest[T, U] {
	return newDecomposeSliceRequest[T, U](input, NewDecomposeOptions())
}

func (r DecomposeSliceRequest[T, U]) Run() ([]U, error) {
	return DecomposeToSlice[T, U](r.input, r.opts)
}

// EnrichRequest is a fluent builder for Enrich.
type EnrichRequest[T any, U any] struct {
	commonRequest[EnrichRequest[T, U], EnrichOptions]
	input T
}

func newEnrichRequest[T any, U any](input T, opts EnrichOptions) EnrichRequest[T, U] {
	return EnrichRequest[T, U]{
		commonRequest: commonRequest[EnrichRequest[T, U], EnrichOptions]{
			opts: opts,
			lift: func(next EnrichOptions) EnrichRequest[T, U] {
				return newEnrichRequest[T, U](input, next)
			},
			mutate: func(current EnrichOptions, fn func(CommonOptions) CommonOptions) EnrichOptions {
				current.CommonOptions = fn(current.CommonOptions)
				return current
			},
		},
		input: input,
	}
}

// Enriching starts a fluent Enrich request.
func Enriching[T any, U any](input T) EnrichRequest[T, U] {
	return newEnrichRequest[T, U](input, NewEnrichOptions())
}

func (r EnrichRequest[T, U]) Run() (EnrichResult[U], error) {
	return Enrich[T, U](r.input, r.opts)
}

// EnrichInPlaceRequest is a fluent builder for EnrichInPlace.
type EnrichInPlaceRequest[T any] struct {
	commonRequest[EnrichInPlaceRequest[T], EnrichOptions]
	input T
}

func newEnrichInPlaceRequest[T any](input T, opts EnrichOptions) EnrichInPlaceRequest[T] {
	return EnrichInPlaceRequest[T]{
		commonRequest: commonRequest[EnrichInPlaceRequest[T], EnrichOptions]{
			opts: opts,
			lift: func(next EnrichOptions) EnrichInPlaceRequest[T] {
				return newEnrichInPlaceRequest(input, next)
			},
			mutate: func(current EnrichOptions, fn func(CommonOptions) CommonOptions) EnrichOptions {
				current.CommonOptions = fn(current.CommonOptions)
				return current
			},
		},
		input: input,
	}
}

// EnrichingInPlace starts a fluent EnrichInPlace request.
func EnrichingInPlace[T any](input T) EnrichInPlaceRequest[T] {
	return newEnrichInPlaceRequest(input, NewEnrichOptions())
}

func (r EnrichInPlaceRequest[T]) Run() (T, error) {
	return EnrichInPlace[T](r.input, r.opts)
}

// NormalizeRequest is a fluent builder for Normalize.
type NormalizeRequest[T any] struct {
	commonRequest[NormalizeRequest[T], NormalizeOptions]
	input T
}

func newNormalizeRequest[T any](input T, opts NormalizeOptions) NormalizeRequest[T] {
	return NormalizeRequest[T]{
		commonRequest: commonRequest[NormalizeRequest[T], NormalizeOptions]{
			opts: opts,
			lift: func(next NormalizeOptions) NormalizeRequest[T] {
				return newNormalizeRequest(input, next)
			},
			mutate: func(current NormalizeOptions, fn func(CommonOptions) CommonOptions) NormalizeOptions {
				current.CommonOptions = fn(current.CommonOptions)
				return current
			},
		},
		input: input,
	}
}

// Normalizing starts a fluent Normalize request.
func Normalizing[T any](input T) NormalizeRequest[T] {
	return newNormalizeRequest(input, NewNormalizeOptions())
}

func (r NormalizeRequest[T]) Run() (NormalizeResult[T], error) {
	return Normalize[T](r.input, r.opts)
}

// NormalizeTextRequest is a fluent builder for NormalizeText.
type NormalizeTextRequest struct {
	commonRequest[NormalizeTextRequest, NormalizeOptions]
	input string
}

func newNormalizeTextRequest(input string, opts NormalizeOptions) NormalizeTextRequest {
	return NormalizeTextRequest{
		commonRequest: commonRequest[NormalizeTextRequest, NormalizeOptions]{
			opts: opts,
			lift: func(next NormalizeOptions) NormalizeTextRequest {
				return newNormalizeTextRequest(input, next)
			},
			mutate: func(current NormalizeOptions, fn func(CommonOptions) CommonOptions) NormalizeOptions {
				current.CommonOptions = fn(current.CommonOptions)
				return current
			},
		},
		input: input,
	}
}

// NormalizingText starts a fluent NormalizeText request.
func NormalizingText(input string) NormalizeTextRequest {
	return newNormalizeTextRequest(input, NewNormalizeOptions())
}

func (r NormalizeTextRequest) Run() (string, error) {
	return NormalizeText(r.input, r.opts)
}

// NormalizeBatchRequest is a fluent builder for NormalizeBatch.
type NormalizeBatchRequest[T any] struct {
	commonRequest[NormalizeBatchRequest[T], NormalizeOptions]
	items []T
}

func newNormalizeBatchRequest[T any](items []T, opts NormalizeOptions) NormalizeBatchRequest[T] {
	return NormalizeBatchRequest[T]{
		commonRequest: commonRequest[NormalizeBatchRequest[T], NormalizeOptions]{
			opts: opts,
			lift: func(next NormalizeOptions) NormalizeBatchRequest[T] {
				return newNormalizeBatchRequest(items, next)
			},
			mutate: func(current NormalizeOptions, fn func(CommonOptions) CommonOptions) NormalizeOptions {
				current.CommonOptions = fn(current.CommonOptions)
				return current
			},
		},
		items: items,
	}
}

// NormalizingBatch starts a fluent NormalizeBatch request.
func NormalizingBatch[T any](items []T) NormalizeBatchRequest[T] {
	return newNormalizeBatchRequest(items, NewNormalizeOptions())
}

func (r NormalizeBatchRequest[T]) Run() ([]NormalizeResult[T], error) {
	return NormalizeBatch[T](r.items, r.opts)
}

// SemanticMatchRequest is a fluent builder for SemanticMatch.
type SemanticMatchRequest[S any, T any] struct {
	commonRequest[SemanticMatchRequest[S, T], MatchOptions]
	sources []S
	targets []T
}

func newSemanticMatchRequest[S any, T any](sources []S, targets []T, opts MatchOptions) SemanticMatchRequest[S, T] {
	return SemanticMatchRequest[S, T]{
		commonRequest: commonRequest[SemanticMatchRequest[S, T], MatchOptions]{
			opts: opts,
			lift: func(next MatchOptions) SemanticMatchRequest[S, T] {
				return newSemanticMatchRequest(sources, targets, next)
			},
			mutate: func(current MatchOptions, fn func(CommonOptions) CommonOptions) MatchOptions {
				current.CommonOptions = fn(current.CommonOptions)
				return current
			},
		},
		sources: sources,
		targets: targets,
	}
}

// Matching starts a fluent SemanticMatch request.
func Matching[S any, T any](sources []S, targets []T) SemanticMatchRequest[S, T] {
	return newSemanticMatchRequest(sources, targets, NewMatchOptions())
}

func (r SemanticMatchRequest[S, T]) By(criteria string) SemanticMatchRequest[S, T] {
	opts := r.opts
	opts.MatchCriteria = criteria
	return r.WithOptions(opts)
}

func (r SemanticMatchRequest[S, T]) Strategy(strategy string) SemanticMatchRequest[S, T] {
	opts := r.opts
	opts.Strategy = strategy
	return r.WithOptions(opts)
}

func (r SemanticMatchRequest[S, T]) Run() (MatchResult[S, T], error) {
	return SemanticMatch[S, T](r.sources, r.targets, r.opts)
}

// MatchOneRequest is a fluent builder for MatchOne.
type MatchOneRequest[S any, T any] struct {
	commonRequest[MatchOneRequest[S, T], MatchOptions]
	source  S
	targets []T
}

func newMatchOneRequest[S any, T any](source S, targets []T, opts MatchOptions) MatchOneRequest[S, T] {
	return MatchOneRequest[S, T]{
		commonRequest: commonRequest[MatchOneRequest[S, T], MatchOptions]{
			opts: opts,
			lift: func(next MatchOptions) MatchOneRequest[S, T] {
				return newMatchOneRequest(source, targets, next)
			},
			mutate: func(current MatchOptions, fn func(CommonOptions) CommonOptions) MatchOptions {
				current.CommonOptions = fn(current.CommonOptions)
				return current
			},
		},
		source:  source,
		targets: targets,
	}
}

// MatchingOne starts a fluent MatchOne request.
func MatchingOne[S any, T any](source S, targets []T) MatchOneRequest[S, T] {
	return newMatchOneRequest(source, targets, NewMatchOptions())
}

func (r MatchOneRequest[S, T]) By(criteria string) MatchOneRequest[S, T] {
	opts := r.opts
	opts.MatchCriteria = criteria
	return r.WithOptions(opts)
}

func (r MatchOneRequest[S, T]) Strategy(strategy string) MatchOneRequest[S, T] {
	opts := r.opts
	opts.Strategy = strategy
	return r.WithOptions(opts)
}

func (r MatchOneRequest[S, T]) Run() ([]MatchPair[S, T], error) {
	return MatchOne[S, T](r.source, r.targets, r.opts)
}

// CritiqueRequest is a fluent builder for Critique.
type CritiqueRequest[T any] struct {
	commonRequest[CritiqueRequest[T], CritiqueOptions]
	input T
}

func newCritiqueRequest[T any](input T, opts CritiqueOptions) CritiqueRequest[T] {
	return CritiqueRequest[T]{
		commonRequest: commonRequest[CritiqueRequest[T], CritiqueOptions]{
			opts: opts,
			lift: func(next CritiqueOptions) CritiqueRequest[T] {
				return newCritiqueRequest(input, next)
			},
			mutate: func(current CritiqueOptions, fn func(CommonOptions) CommonOptions) CritiqueOptions {
				current.CommonOptions = fn(current.CommonOptions)
				return current
			},
		},
		input: input,
	}
}

// Critiquing starts a fluent Critique request.
func Critiquing[T any](input T) CritiqueRequest[T] {
	return newCritiqueRequest(input, NewCritiqueOptions())
}

func (r CritiqueRequest[T]) Run() (CritiqueResult, error) {
	return Critique[T](r.input, r.opts)
}

// SynthesizeRequest is a fluent builder for Synthesize.
type SynthesizeRequest[T any] struct {
	commonRequest[SynthesizeRequest[T], SynthesizeOptions]
	sources []any
}

func newSynthesizeRequest[T any](sources []any, opts SynthesizeOptions) SynthesizeRequest[T] {
	return SynthesizeRequest[T]{
		commonRequest: commonRequest[SynthesizeRequest[T], SynthesizeOptions]{
			opts: opts,
			lift: func(next SynthesizeOptions) SynthesizeRequest[T] {
				return newSynthesizeRequest[T](sources, next)
			},
			mutate: func(current SynthesizeOptions, fn func(CommonOptions) CommonOptions) SynthesizeOptions {
				current.CommonOptions = fn(current.CommonOptions)
				return current
			},
		},
		sources: sources,
	}
}

// Synthesizing starts a fluent Synthesize request.
func Synthesizing[T any](sources []any) SynthesizeRequest[T] {
	return newSynthesizeRequest[T](sources, NewSynthesizeOptions())
}

func (r SynthesizeRequest[T]) Strategy(strategy string) SynthesizeRequest[T] {
	opts := r.opts
	opts.Strategy = strategy
	return r.WithOptions(opts)
}

func (r SynthesizeRequest[T]) Run() (SynthesizeResult[T], error) {
	return Synthesize[T](r.sources, r.opts)
}

// PredictRequest is a fluent builder for Predict.
type PredictRequest[T any] struct {
	commonRequest[PredictRequest[T], PredictOptions]
	input any
}

func newPredictRequest[T any](input any, opts PredictOptions) PredictRequest[T] {
	return PredictRequest[T]{
		commonRequest: commonRequest[PredictRequest[T], PredictOptions]{
			opts: opts,
			lift: func(next PredictOptions) PredictRequest[T] {
				return newPredictRequest[T](input, next)
			},
			mutate: func(current PredictOptions, fn func(CommonOptions) CommonOptions) PredictOptions {
				current.CommonOptions = fn(current.CommonOptions)
				return current
			},
		},
		input: input,
	}
}

// Predicting starts a fluent Predict request.
func Predicting[T any](input any) PredictRequest[T] {
	return newPredictRequest[T](input, NewPredictOptions())
}

func (r PredictRequest[T]) Horizon(horizon string) PredictRequest[T] {
	opts := r.opts
	opts.Horizon = horizon
	return r.WithOptions(opts)
}

func (r PredictRequest[T]) Run() (PredictResult[T], error) {
	return Predict[T](r.input, r.opts)
}

// VerifyRequest is a fluent builder for Verify.
type VerifyRequest struct {
	commonRequest[VerifyRequest, VerifyOptions]
	input string
}

func newVerifyRequest(input string, opts VerifyOptions) VerifyRequest {
	return VerifyRequest{
		commonRequest: commonRequest[VerifyRequest, VerifyOptions]{
			opts: opts,
			lift: func(next VerifyOptions) VerifyRequest {
				return newVerifyRequest(input, next)
			},
			mutate: func(current VerifyOptions, fn func(CommonOptions) CommonOptions) VerifyOptions {
				current.CommonOptions = fn(current.CommonOptions)
				return current
			},
		},
		input: input,
	}
}

// Verifying starts a fluent Verify request.
func Verifying(input string) VerifyRequest {
	return newVerifyRequest(input, NewVerifyOptions())
}

func (r VerifyRequest) Run() (VerifyResult, error) {
	return Verify(r.input, r.opts)
}

// VerifyClaimRequest is a fluent builder for VerifyClaim.
type VerifyClaimRequest struct {
	commonRequest[VerifyClaimRequest, VerifyOptions]
	claim string
}

func newVerifyClaimRequest(claim string, opts VerifyOptions) VerifyClaimRequest {
	return VerifyClaimRequest{
		commonRequest: commonRequest[VerifyClaimRequest, VerifyOptions]{
			opts: opts,
			lift: func(next VerifyOptions) VerifyClaimRequest {
				return newVerifyClaimRequest(claim, next)
			},
			mutate: func(current VerifyOptions, fn func(CommonOptions) CommonOptions) VerifyOptions {
				current.CommonOptions = fn(current.CommonOptions)
				return current
			},
		},
		claim: claim,
	}
}

// VerifyingClaim starts a fluent VerifyClaim request.
func VerifyingClaim(claim string) VerifyClaimRequest {
	return newVerifyClaimRequest(claim, NewVerifyOptions())
}

func (r VerifyClaimRequest) Run() (ClaimVerification, error) {
	return VerifyClaim(r.claim, r.opts)
}
