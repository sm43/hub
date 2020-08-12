// Code generated by goa v3.2.2, DO NOT EDIT.
//
// resource HTTP server encoders and decoders
//
// Command:
// $ goa gen github.com/tektoncd/hub/api/design

package server

import (
	"context"
	"net/http"
	"strconv"

	resourceviews "github.com/tektoncd/hub/api/gen/resource/views"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// EncodeQueryResponse returns an encoder for responses returned by the
// resource Query endpoint.
func EncodeQueryResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res := v.(resourceviews.ResourceCollection)
		enc := encoder(ctx, w)
		body := NewResourceResponseWithoutVersionCollection(res.Projected)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeQueryRequest returns a decoder for requests sent to the resource Query
// endpoint.
func DecodeQueryRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			name  string
			type_ string
			limit uint
			err   error
		)
		nameRaw := r.URL.Query().Get("name")
		if nameRaw != "" {
			name = nameRaw
		}
		type_Raw := r.URL.Query().Get("type")
		if type_Raw != "" {
			type_ = type_Raw
		}
		if !(type_ == "task" || type_ == "pipeline" || type_ == "") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("type_", type_, []interface{}{"task", "pipeline", ""}))
		}
		{
			limitRaw := r.URL.Query().Get("limit")
			if limitRaw == "" {
				limit = 100
			} else {
				v, err2 := strconv.ParseUint(limitRaw, 10, strconv.IntSize)
				if err2 != nil {
					err = goa.MergeErrors(err, goa.InvalidFieldTypeError("limit", limitRaw, "unsigned integer"))
				}
				limit = uint(v)
			}
		}
		if err != nil {
			return nil, err
		}
		payload := NewQueryPayload(name, type_, limit)

		return payload, nil
	}
}

// EncodeQueryError returns an encoder for errors returned by the Query
// resource endpoint.
func EncodeQueryError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		en, ok := v.(ErrorNamer)
		if !ok {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "internal-error":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewQueryInternalErrorResponseBody(res)
			}
			w.Header().Set("goa-error", "internal-error")
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		case "not-found":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewQueryNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", "not-found")
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeListResponse returns an encoder for responses returned by the resource
// List endpoint.
func EncodeListResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res := v.(resourceviews.ResourceCollection)
		enc := encoder(ctx, w)
		body := NewResourceResponseWithoutVersionCollection(res.Projected)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeListRequest returns a decoder for requests sent to the resource List
// endpoint.
func DecodeListRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			limit uint
			err   error
		)
		{
			limitRaw := r.URL.Query().Get("limit")
			if limitRaw == "" {
				limit = 100
			} else {
				v, err2 := strconv.ParseUint(limitRaw, 10, strconv.IntSize)
				if err2 != nil {
					err = goa.MergeErrors(err, goa.InvalidFieldTypeError("limit", limitRaw, "unsigned integer"))
				}
				limit = uint(v)
			}
		}
		if err != nil {
			return nil, err
		}
		payload := NewListPayload(limit)

		return payload, nil
	}
}

// EncodeListError returns an encoder for errors returned by the List resource
// endpoint.
func EncodeListError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		en, ok := v.(ErrorNamer)
		if !ok {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "internal-error":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewListInternalErrorResponseBody(res)
			}
			w.Header().Set("goa-error", "internal-error")
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeVersionsByIDResponse returns an encoder for responses returned by the
// resource VersionsByID endpoint.
func EncodeVersionsByIDResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res := v.(*resourceviews.Versions)
		enc := encoder(ctx, w)
		body := NewVersionsByIDResponseBody(res.Projected)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeVersionsByIDRequest returns a decoder for requests sent to the
// resource VersionsByID endpoint.
func DecodeVersionsByIDRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			id  uint
			err error

			params = mux.Vars(r)
		)
		{
			idRaw := params["id"]
			v, err2 := strconv.ParseUint(idRaw, 10, strconv.IntSize)
			if err2 != nil {
				err = goa.MergeErrors(err, goa.InvalidFieldTypeError("id", idRaw, "unsigned integer"))
			}
			id = uint(v)
		}
		if err != nil {
			return nil, err
		}
		payload := NewVersionsByIDPayload(id)

		return payload, nil
	}
}

// EncodeVersionsByIDError returns an encoder for errors returned by the
// VersionsByID resource endpoint.
func EncodeVersionsByIDError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		en, ok := v.(ErrorNamer)
		if !ok {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "internal-error":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewVersionsByIDInternalErrorResponseBody(res)
			}
			w.Header().Set("goa-error", "internal-error")
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		case "not-found":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewVersionsByIDNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", "not-found")
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeByTypeNameVersionResponse returns an encoder for responses returned by
// the resource ByTypeNameVersion endpoint.
func EncodeByTypeNameVersionResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res := v.(*resourceviews.Version)
		enc := encoder(ctx, w)
		body := NewByTypeNameVersionResponseBody(res.Projected)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeByTypeNameVersionRequest returns a decoder for requests sent to the
// resource ByTypeNameVersion endpoint.
func DecodeByTypeNameVersionRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			type_   string
			name    string
			version string
			err     error

			params = mux.Vars(r)
		)
		type_ = params["type"]
		if !(type_ == "task" || type_ == "pipeline") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("type_", type_, []interface{}{"task", "pipeline"}))
		}
		name = params["name"]
		version = params["version"]
		if err != nil {
			return nil, err
		}
		payload := NewByTypeNameVersionPayload(type_, name, version)

		return payload, nil
	}
}

// EncodeByTypeNameVersionError returns an encoder for errors returned by the
// ByTypeNameVersion resource endpoint.
func EncodeByTypeNameVersionError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		en, ok := v.(ErrorNamer)
		if !ok {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "internal-error":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewByTypeNameVersionInternalErrorResponseBody(res)
			}
			w.Header().Set("goa-error", "internal-error")
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		case "not-found":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewByTypeNameVersionNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", "not-found")
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeByVersionIDResponse returns an encoder for responses returned by the
// resource ByVersionId endpoint.
func EncodeByVersionIDResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res := v.(*resourceviews.Version)
		enc := encoder(ctx, w)
		body := NewByVersionIDResponseBody(res.Projected)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeByVersionIDRequest returns a decoder for requests sent to the resource
// ByVersionId endpoint.
func DecodeByVersionIDRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			versionID uint
			err       error

			params = mux.Vars(r)
		)
		{
			versionIDRaw := params["versionID"]
			v, err2 := strconv.ParseUint(versionIDRaw, 10, strconv.IntSize)
			if err2 != nil {
				err = goa.MergeErrors(err, goa.InvalidFieldTypeError("versionID", versionIDRaw, "unsigned integer"))
			}
			versionID = uint(v)
		}
		if err != nil {
			return nil, err
		}
		payload := NewByVersionIDPayload(versionID)

		return payload, nil
	}
}

// EncodeByVersionIDError returns an encoder for errors returned by the
// ByVersionId resource endpoint.
func EncodeByVersionIDError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		en, ok := v.(ErrorNamer)
		if !ok {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "internal-error":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewByVersionIDInternalErrorResponseBody(res)
			}
			w.Header().Set("goa-error", "internal-error")
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		case "not-found":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewByVersionIDNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", "not-found")
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeByTypeNameResponse returns an encoder for responses returned by the
// resource ByTypeName endpoint.
func EncodeByTypeNameResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res := v.(resourceviews.ResourceCollection)
		enc := encoder(ctx, w)
		body := NewResourceResponseWithoutVersionCollection(res.Projected)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeByTypeNameRequest returns a decoder for requests sent to the resource
// ByTypeName endpoint.
func DecodeByTypeNameRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			type_ string
			name  string
			err   error

			params = mux.Vars(r)
		)
		type_ = params["type"]
		if !(type_ == "task" || type_ == "pipeline") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("type_", type_, []interface{}{"task", "pipeline"}))
		}
		name = params["name"]
		if err != nil {
			return nil, err
		}
		payload := NewByTypeNamePayload(type_, name)

		return payload, nil
	}
}

// EncodeByTypeNameError returns an encoder for errors returned by the
// ByTypeName resource endpoint.
func EncodeByTypeNameError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		en, ok := v.(ErrorNamer)
		if !ok {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "internal-error":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewByTypeNameInternalErrorResponseBody(res)
			}
			w.Header().Set("goa-error", "internal-error")
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		case "not-found":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewByTypeNameNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", "not-found")
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeByIDResponse returns an encoder for responses returned by the resource
// ById endpoint.
func EncodeByIDResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res := v.(*resourceviews.Resource)
		enc := encoder(ctx, w)
		body := NewByIDResponseBody(res.Projected)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeByIDRequest returns a decoder for requests sent to the resource ById
// endpoint.
func DecodeByIDRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			id  uint
			err error

			params = mux.Vars(r)
		)
		{
			idRaw := params["id"]
			v, err2 := strconv.ParseUint(idRaw, 10, strconv.IntSize)
			if err2 != nil {
				err = goa.MergeErrors(err, goa.InvalidFieldTypeError("id", idRaw, "unsigned integer"))
			}
			id = uint(v)
		}
		if err != nil {
			return nil, err
		}
		payload := NewByIDPayload(id)

		return payload, nil
	}
}

// EncodeByIDError returns an encoder for errors returned by the ById resource
// endpoint.
func EncodeByIDError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		en, ok := v.(ErrorNamer)
		if !ok {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "internal-error":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewByIDInternalErrorResponseBody(res)
			}
			w.Header().Set("goa-error", "internal-error")
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		case "not-found":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewByIDNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", "not-found")
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// marshalResourceviewsResourceViewToResourceResponseWithoutVersion builds a
// value of type *ResourceResponseWithoutVersion from a value of type
// *resourceviews.ResourceView.
func marshalResourceviewsResourceViewToResourceResponseWithoutVersion(v *resourceviews.ResourceView) *ResourceResponseWithoutVersion {
	res := &ResourceResponseWithoutVersion{
		ID:     *v.ID,
		Name:   *v.Name,
		Type:   *v.Type,
		Rating: *v.Rating,
	}
	if v.Catalog != nil {
		res.Catalog = marshalResourceviewsCatalogViewToCatalogResponse(v.Catalog)
	}
	if v.LatestVersion != nil {
		res.LatestVersion = marshalResourceviewsVersionViewToVersionResponseWithoutResource(v.LatestVersion)
	}
	if v.Tags != nil {
		res.Tags = make([]*TagResponse, len(v.Tags))
		for i, val := range v.Tags {
			res.Tags[i] = marshalResourceviewsTagViewToTagResponse(val)
		}
	}

	return res
}

// marshalResourceviewsCatalogViewToCatalogResponse builds a value of type
// *CatalogResponse from a value of type *resourceviews.CatalogView.
func marshalResourceviewsCatalogViewToCatalogResponse(v *resourceviews.CatalogView) *CatalogResponse {
	res := &CatalogResponse{
		ID:   *v.ID,
		Type: *v.Type,
	}

	return res
}

// marshalResourceviewsVersionViewToVersionResponseWithoutResource builds a
// value of type *VersionResponseWithoutResource from a value of type
// *resourceviews.VersionView.
func marshalResourceviewsVersionViewToVersionResponseWithoutResource(v *resourceviews.VersionView) *VersionResponseWithoutResource {
	res := &VersionResponseWithoutResource{
		ID:                  *v.ID,
		Version:             *v.Version,
		DisplayName:         *v.DisplayName,
		Description:         *v.Description,
		MinPipelinesVersion: *v.MinPipelinesVersion,
		RawURL:              *v.RawURL,
		WebURL:              *v.WebURL,
		UpdatedAt:           *v.UpdatedAt,
	}

	return res
}

// marshalResourceviewsTagViewToTagResponse builds a value of type *TagResponse
// from a value of type *resourceviews.TagView.
func marshalResourceviewsTagViewToTagResponse(v *resourceviews.TagView) *TagResponse {
	res := &TagResponse{
		ID:   *v.ID,
		Name: *v.Name,
	}

	return res
}

// marshalResourceviewsVersionViewToVersionResponseBodyMin builds a value of
// type *VersionResponseBodyMin from a value of type *resourceviews.VersionView.
func marshalResourceviewsVersionViewToVersionResponseBodyMin(v *resourceviews.VersionView) *VersionResponseBodyMin {
	res := &VersionResponseBodyMin{
		ID:      *v.ID,
		Version: *v.Version,
		RawURL:  *v.RawURL,
		WebURL:  *v.WebURL,
	}

	return res
}

// marshalResourceviewsResourceViewToResourceResponseBodyInfo builds a value of
// type *ResourceResponseBodyInfo from a value of type
// *resourceviews.ResourceView.
func marshalResourceviewsResourceViewToResourceResponseBodyInfo(v *resourceviews.ResourceView) *ResourceResponseBodyInfo {
	res := &ResourceResponseBodyInfo{
		ID:     *v.ID,
		Name:   *v.Name,
		Type:   *v.Type,
		Rating: *v.Rating,
	}
	if v.Catalog != nil {
		res.Catalog = marshalResourceviewsCatalogViewToCatalogResponseBody(v.Catalog)
	}
	if v.Tags != nil {
		res.Tags = make([]*TagResponseBody, len(v.Tags))
		for i, val := range v.Tags {
			res.Tags[i] = marshalResourceviewsTagViewToTagResponseBody(val)
		}
	}

	return res
}

// marshalResourceviewsCatalogViewToCatalogResponseBody builds a value of type
// *CatalogResponseBody from a value of type *resourceviews.CatalogView.
func marshalResourceviewsCatalogViewToCatalogResponseBody(v *resourceviews.CatalogView) *CatalogResponseBody {
	res := &CatalogResponseBody{
		ID:   *v.ID,
		Type: *v.Type,
	}

	return res
}

// marshalResourceviewsTagViewToTagResponseBody builds a value of type
// *TagResponseBody from a value of type *resourceviews.TagView.
func marshalResourceviewsTagViewToTagResponseBody(v *resourceviews.TagView) *TagResponseBody {
	res := &TagResponseBody{
		ID:   *v.ID,
		Name: *v.Name,
	}

	return res
}

// marshalResourceviewsVersionViewToVersionResponseBodyWithoutResource builds a
// value of type *VersionResponseBodyWithoutResource from a value of type
// *resourceviews.VersionView.
func marshalResourceviewsVersionViewToVersionResponseBodyWithoutResource(v *resourceviews.VersionView) *VersionResponseBodyWithoutResource {
	res := &VersionResponseBodyWithoutResource{
		ID:                  *v.ID,
		Version:             *v.Version,
		DisplayName:         *v.DisplayName,
		Description:         *v.Description,
		MinPipelinesVersion: *v.MinPipelinesVersion,
		RawURL:              *v.RawURL,
		WebURL:              *v.WebURL,
		UpdatedAt:           *v.UpdatedAt,
	}

	return res
}

// marshalResourceviewsVersionViewToVersionResponseBodyTiny builds a value of
// type *VersionResponseBodyTiny from a value of type
// *resourceviews.VersionView.
func marshalResourceviewsVersionViewToVersionResponseBodyTiny(v *resourceviews.VersionView) *VersionResponseBodyTiny {
	res := &VersionResponseBodyTiny{
		ID:      *v.ID,
		Version: *v.Version,
	}

	return res
}
