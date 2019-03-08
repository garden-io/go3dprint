package svg

/*
#cgo pkg-config: --cflags --libs librsvg-2.0 cairo-pdf
#include <librsvg/rsvg.h>

static void
sizeCallback(int *width, int *height, gpointer data) {
	RsvgDimensionData *size = data;
	*width = size->width;
	*height = size->height;
}

static void
setOutputSize(RsvgHandle *handle, RsvgDimensionData *size) {
	// TODO do not use this deprecated function - use
	// cairo matrix instead, as suggested in:
	// https://developer.gnome.org/rsvg/stable/RsvgHandle.html#rsvg-handle-set-size-callback
	rsvg_handle_set_size_callback(handle, sizeCallback, size, 0);
}
*/
import "C"
import (
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"reflect"
	"unsafe"
)

// Render reads an SVG from the given reader and renders
// it into an image of the given size.
func Render(svg io.Reader, size image.Point) (*image.RGBA, error) {
	// TODO use GInputStream
	svgData, err := ioutil.ReadAll(svg)
	if err != nil {
		return nil, err
	}
	if len(svgData) == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	var gerr *C.GError
	handle := C.rsvg_handle_new_from_data((*C.guint8)(unsafe.Pointer(&svgData[0])), C.gsize(len(svgData)), &gerr)
	if gerr != nil {
		return nil, fmt.Errorf("cannot make new handle: %s", C.GoString((*C.char)(unsafe.Pointer(gerr.message))))
	}
	dimSize := C.RsvgDimensionData{
		width:  C.int(size.X),
		height: C.int(size.Y),
	}
	C.setOutputSize(handle, &dimSize)

	surface := C.cairo_image_surface_create(C.CAIRO_FORMAT_ARGB32, C.int(size.X), C.int(size.Y))
	cr := C.cairo_create(surface)
	C.rsvg_handle_render_cairo(handle, cr)
	if status := C.cairo_status(cr); status != 0 {
		errStr := C.cairo_status_to_string(status)
		return nil, fmt.Errorf("cannot render svg: %s", errStr)
	}
	C.cairo_surface_flush(surface)

	numPixels := size.X * size.Y
	dataBytes := C.cairo_image_surface_get_data(surface)
	dataHdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(dataBytes)),
		Len:  numPixels,
		Cap:  numPixels,
	}
	data := *(*[]uint32)(unsafe.Pointer(&dataHdr))

	goData := make([]uint8, numPixels*4)

	for i, pix := range data {
		d := goData[i*4:]
		d[0] = uint8(pix >> 16) // R
		d[1] = uint8(pix >> 8)  // G
		d[2] = uint8(pix >> 0)  // B
		d[3] = uint8(pix >> 24) // A
	}
	return &image.RGBA{
		Pix:    goData,
		Stride: size.X * 4,
		Rect: image.Rectangle{
			Max: size,
		},
	}, nil
}
