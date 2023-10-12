#include <dlib/dnn.h>
#include <dlib/image_loader/image_loader.h>
#include <dlib/graph_utils.h>
#include "jpeg_mem_loader.h"
#include "object_detector.h"

using namespace dlib;

typedef scan_fhog_pyramid<pyramid_down<6> > image_scanner_type;

static const size_t RECT_LEN = 4;
static const size_t DESCR_LEN = 128;
static const size_t SHAPE_LEN = 2;
static const size_t RECT_SIZE = RECT_LEN * sizeof(long);
static const size_t DESCR_SIZE = DESCR_LEN * sizeof(float);
static const size_t SHAPE_SIZE = SHAPE_LEN * sizeof(long);

class ObjRec
{
public:
	ObjRec(const char *model)
	{
		object_detector<image_scanner_type> detector;
		deserialize(model) >> detector;
		detector_ = detector;
	}
	std::vector<rectangle> Recognize(const matrix<rgb_pixel> &img)
	{
		std::vector<rectangle> rects;

		std::lock_guard<std::mutex> lock(detector_mutex_);
		rects = detector_(img);

		if (rects.size() == 0)
		{
			return rects;
		}

		std::sort(rects.begin(), rects.end());

		return rects;
	}

private:
	std::mutex detector_mutex_;
	object_detector<image_scanner_type> detector_;
};

objrec *objrec_init(const char *model_dir)
{
	objrec *rec = (objrec *)calloc(1, sizeof(objrec));
	try
	{
		ObjRec *cls = new ObjRec(model_dir);
		rec->cls = (void *)cls;
	}
	catch (serialization_error &e)
	{
		rec->err_str = strdup(e.what());
		rec->err_code = SERIALIZATION_ERROR;
	}
	catch (std::exception &e)
	{
		rec->err_str = strdup(e.what());
		rec->err_code = UNKNOWN_ERROR;
	}
	return rec;
}

objret *objrec_recognize(objrec *rec, const uint8_t *img_data, int len)
{
	objret *ret = (objret *)calloc(1, sizeof(objret));
	ObjRec *cls = (ObjRec *)(rec->cls);
	matrix<rgb_pixel> img;
	std::vector<rectangle> rects;
	try
	{
		load_mem_jpeg(img, img_data, len);
		rects = cls->Recognize(img);
	}
	catch (image_load_error &e)
	{
		ret->err_str = strdup(e.what());
		ret->err_code = IMAGE_LOAD_ERROR;
		return ret;
	}
	catch (std::exception &e)
	{
		ret->err_str = strdup(e.what());
		ret->err_code = UNKNOWN_ERROR;
		return ret;
	}

	ret->rectCount = rects.size();
	if (ret->rectCount == 0)
	{
		return ret;
	}

	ret->rectangles = (long *)malloc(ret->rectCount * RECT_SIZE);
	for (int i = 0; i < ret->rectCount; i++)
	{
		long *dst = ret->rectangles + i * RECT_LEN;
		dst[0] = rects[i].left();
		dst[1] = rects[i].top();
		dst[2] = rects[i].right();
		dst[3] = rects[i].bottom();
	}

	return ret;
}

void objrec_free(objrec *rec)
{
	if (rec)
	{
		if (rec->cls)
		{
			ObjRec *cls = (ObjRec *)(rec->cls);
			delete cls;
			rec->cls = NULL;
		}
		free(rec);
	}
}