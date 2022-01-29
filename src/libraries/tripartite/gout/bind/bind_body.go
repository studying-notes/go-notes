package main

import (
	"encoding/base64"
	"fmt"
	"github.com/guonaihong/gout"
	"io/ioutil"
	"time"
)

type ImageItem struct {
	ImageData struct {
		Image string `json:"image"`
	} `json:"ImageData"`
	Light   int `json:"light"`
	PageIdx int `json:"page_idx"`
}

type Payload struct {
	ProcessParam struct {
		DoublePageSpread bool   `json:"doublePageSpread"`
		Scenario         string `json:"scenario"`
	} `json:"processParam"`
	List []ImageItem `json:"List"`
}

func NewPayload(image string) (payload Payload) {
	payload.ProcessParam.DoublePageSpread = true
	payload.ProcessParam.Scenario = "FullProcess"

	imageItem := ImageItem{Light: 6, PageIdx: 0}
	imageItem.ImageData.Image = image

	payload.List = []ImageItem{imageItem}

	return payload
}

type Response struct {
	ChipPage      int `json:"ChipPage"`
	ContainerList struct {
		Count int `json:"Count"`
		List  []struct {
			ResultType int `json:"result_type"`
			BufLength  int `json:"buf_length"`
			Light      int `json:"light"`
			ListIdx    int `json:"list_idx"`
			PageIdx    int `json:"page_idx"`
			Status     struct {
				DetailsOptical struct {
					DocType       int `json:"docType"`
					Expiry        int `json:"expiry"`
					ImageQA       int `json:"imageQA"`
					Mrz           int `json:"mrz"`
					OverallStatus int `json:"overallStatus"`
					PagesCount    int `json:"pagesCount"`
					Security      int `json:"security"`
					Text          int `json:"text"`
				} `json:"detailsOptical"`
				DetailsRFID struct {
					AA            int `json:"AA"`
					BAC           int `json:"BAC"`
					CA            int `json:"CA"`
					PA            int `json:"PA"`
					PACE          int `json:"PACE"`
					TA            int `json:"TA"`
					OverallStatus int `json:"overallStatus"`
				} `json:"detailsRFID"`
				Optical       int `json:"optical"`
				OverallStatus int `json:"overallStatus"`
				Portrait      int `json:"portrait"`
				Rfid          int `json:"rfid"`
				StopList      int `json:"stopList"`
			} `json:"Status,omitempty"`
			Images struct {
				AvailableSourceList []struct {
					ContainerType int    `json:"containerType"`
					Source        string `json:"source"`
				} `json:"availableSourceList"`
				FieldList []struct {
					FieldName string `json:"fieldName"`
					FieldType int    `json:"fieldType"`
					ValueList []struct {
						ContainerType int `json:"containerType"`
						FieldRect     struct {
							Bottom int `json:"bottom"`
							Left   int `json:"left"`
							Right  int `json:"right"`
							Top    int `json:"top"`
						} `json:"fieldRect,omitempty"`
						LightIndex        int    `json:"lightIndex"`
						OriginalPageIndex int    `json:"originalPageIndex"`
						PageIndex         int    `json:"pageIndex"`
						Source            string `json:"source"`
						Value             string `json:"value"`
					} `json:"valueList"`
				} `json:"fieldList"`
			} `json:"Images,omitempty"`
			Text struct {
				AvailableSourceList []struct {
					ContainerType  int    `json:"containerType"`
					Source         string `json:"source"`
					ValidityStatus int    `json:"validityStatus"`
				} `json:"availableSourceList"`
				ComparisonStatus int    `json:"comparisonStatus"`
				DateFormat       string `json:"dateFormat"`
				FieldList        []struct {
					ComparisonList   []interface{} `json:"comparisonList"`
					ComparisonStatus int           `json:"comparisonStatus"`
					FieldName        string        `json:"fieldName"`
					FieldType        int           `json:"fieldType"`
					Lcid             int           `json:"lcid"`
					LcidName         string        `json:"lcidName"`
					Status           int           `json:"status"`
					ValidityList     []struct {
						Source string `json:"source"`
						Status int    `json:"status"`
					} `json:"validityList"`
					ValidityStatus int    `json:"validityStatus"`
					Value          string `json:"value"`
					ValueList      []struct {
						ContainerType int `json:"containerType"`
						FieldRect     struct {
							Bottom int `json:"bottom"`
							Left   int `json:"left"`
							Right  int `json:"right"`
							Top    int `json:"top"`
						} `json:"fieldRect,omitempty"`
						OriginalSymbols []struct {
							Code        int `json:"code"`
							Probability int `json:"probability"`
							Rect        struct {
								Bottom int `json:"bottom"`
								Left   int `json:"left"`
								Right  int `json:"right"`
								Top    int `json:"top"`
							} `json:"rect"`
						} `json:"originalSymbols"`
						OriginalValidity int    `json:"originalValidity"`
						OriginalValue    string `json:"originalValue,omitempty"`
						PageIndex        int    `json:"pageIndex"`
						Probability      int    `json:"probability"`
						Source           string `json:"source"`
						Value            string `json:"value"`
					} `json:"valueList"`
				} `json:"fieldList"`
				Status         int `json:"status"`
				ValidityStatus int `json:"validityStatus"`
			} `json:"Text,omitempty"`
			ListVerifiedFields struct {
				Count       int    `json:"Count"`
				PDateFormat string `json:"pDateFormat"`
				PFieldMaps  []struct {
					FieldType   int    `json:"FieldType"`
					FieldVisual string `json:"Field_Visual"`
					Matrix      []int  `json:"Matrix"`
					WFieldType  int    `json:"wFieldType"`
					WLCID       int    `json:"wLCID"`
				} `json:"pFieldMaps"`
			} `json:"ListVerifiedFields,omitempty"`
			DocumentPosition struct {
				Angle  float64 `json:"Angle"`
				Center struct {
					X int `json:"x"`
					Y int `json:"y"`
				} `json:"Center"`
				Dpi        int `json:"Dpi"`
				Height     int `json:"Height"`
				Inverse    int `json:"Inverse"`
				LeftBottom struct {
					X int `json:"x"`
					Y int `json:"y"`
				} `json:"LeftBottom"`
				LeftTop struct {
					X int `json:"x"`
					Y int `json:"y"`
				} `json:"LeftTop"`
				ObjArea        int `json:"ObjArea"`
				ObjIntAngleDev int `json:"ObjIntAngleDev"`
				PerspectiveTr  int `json:"PerspectiveTr"`
				ResultStatus   int `json:"ResultStatus"`
				RightBottom    struct {
					X int `json:"x"`
					Y int `json:"y"`
				} `json:"RightBottom"`
				RightTop struct {
					X int `json:"x"`
					Y int `json:"y"`
				} `json:"RightTop"`
				Width     int `json:"Width"`
				DocFormat int `json:"docFormat"`
			} `json:"DocumentPosition,omitempty"`
			ImageQualityCheckList struct {
				Count int `json:"Count"`
				List  []struct {
					Areas struct {
						Count int `json:"Count"`
						List  []struct {
							Bottom int `json:"bottom"`
							Left   int `json:"left"`
							Right  int `json:"right"`
							Top    int `json:"top"`
						} `json:"List"`
						Points []interface{} `json:"Points"`
					} `json:"areas,omitempty"`
					FeatureType int     `json:"featureType"`
					Mean        float64 `json:"mean"`
					Probability int     `json:"probability"`
					Result      int     `json:"result"`
					StdDev      float64 `json:"std_dev"`
					Type        int     `json:"type"`
				} `json:"List"`
				Result int `json:"result"`
			} `json:"ImageQualityCheckList,omitempty"`
			OneCandidate struct {
				AuthenticityNecessaryLights int    `json:"AuthenticityNecessaryLights"`
				CheckAuthenticity           int    `json:"CheckAuthenticity"`
				DocumentName                string `json:"DocumentName"`
				FDSIDList                   struct {
					Count        int    `json:"Count"`
					ICAOCode     string `json:"ICAOCode"`
					List         []int  `json:"List"`
					DCountryName string `json:"dCountryName"`
					DFormat      int    `json:"dFormat"`
					DMRZ         bool   `json:"dMRZ"`
					DType        int    `json:"dType"`
					DYear        string `json:"dYear"`
				} `json:"FDSIDList"`
				ID              int     `json:"ID"`
				NecessaryLights int     `json:"NecessaryLights"`
				OVIExp          int     `json:"OVIExp"`
				P               float64 `json:"P"`
				RFIDPresence    int     `json:"RFID_Presence"`
				Rotated180      int     `json:"Rotated180"`
				RotationAngle   int     `json:"RotationAngle"`
				UVExp           int     `json:"UVExp"`
			} `json:"OneCandidate,omitempty"`
			DocVisualExtendedInfo struct {
				NFields      int `json:"nFields"`
				PArrayFields []struct {
					BufLength int    `json:"Buf_Length"`
					BufText   string `json:"Buf_Text"`
					FieldMask string `json:"FieldMask"`
					FieldName string `json:"FieldName"`
					FieldRect struct {
						Bottom int `json:"bottom"`
						Left   int `json:"left"`
						Right  int `json:"right"`
						Top    int `json:"top"`
					} `json:"FieldRect"`
					FieldType     int `json:"FieldType"`
					InComparison  int `json:"InComparison"`
					Reserved2     int `json:"Reserved2"`
					Reserved3     int `json:"Reserved3"`
					StringsCount  int `json:"StringsCount"`
					StringsResult []struct {
						Reserved     int `json:"Reserved"`
						StringResult []struct {
							BaseLineBottom   int `json:"BaseLineBottom"`
							BaseLineTop      int `json:"BaseLineTop"`
							CandidatesCount  int `json:"CandidatesCount"`
							ListOfCandidates []struct {
								Class             int `json:"Class"`
								SubClass          int `json:"SubClass"`
								SymbolCode        int `json:"SymbolCode"`
								SymbolProbability int `json:"SymbolProbability"`
							} `json:"ListOfCandidates"`
							Reserved   int `json:"Reserved"`
							SymbolRect struct {
								Bottom int `json:"bottom"`
								Left   int `json:"left"`
								Right  int `json:"right"`
								Top    int `json:"top"`
							} `json:"SymbolRect"`
						} `json:"StringResult"`
						SymbolsCount int `json:"SymbolsCount"`
					} `json:"StringsResult"`
					Validity   int `json:"Validity"`
					WFieldType int `json:"wFieldType"`
					WLCID      int `json:"wLCID"`
				} `json:"pArrayFields"`
			} `json:"DocVisualExtendedInfo,omitempty"`
			DocGraphicsInfo struct {
				NFields      int `json:"nFields"`
				PArrayFields []struct {
					FieldName string `json:"FieldName"`
					FieldRect struct {
						Bottom int `json:"bottom"`
						Left   int `json:"left"`
						Right  int `json:"right"`
						Top    int `json:"top"`
					} `json:"FieldRect"`
					FieldType int `json:"FieldType"`
					Image     struct {
						Format string `json:"format"`
						Image  string `json:"image"`
					} `json:"image"`
				} `json:"pArrayFields"`
			} `json:"DocGraphicsInfo,omitempty"`
			FaceDetection struct {
				Count               int `json:"Count"`
				CountFalseDetection int `json:"CountFalseDetection"`
				Res                 []struct {
					CoincidenceToPhotoArea int `json:"CoincidenceToPhotoArea"`
					LightType              int `json:"LightType"`
					Orientation            int `json:"Orientation"`
					Probability            int `json:"Probability"`
					RectPhoto              struct {
						Bottom int `json:"bottom"`
						Left   int `json:"left"`
						Right  int `json:"right"`
						Top    int `json:"top"`
					} `json:"Rect_Photo"`
					Reserved int `json:"Reserved"`
					PRects   struct {
						Bottom int `json:"bottom"`
						Left   int `json:"left"`
						Right  int `json:"right"`
						Top    int `json:"top"`
					} `json:"pRects"`
				} `json:"Res"`
				Reserved1 int `json:"Reserved1"`
				Reserved2 int `json:"Reserved2"`
			} `json:"FaceDetection,omitempty"`
		} `json:"List"`
	} `json:"ContainerList"`
	CoreLibResultCode  int `json:"CoreLibResultCode"`
	ProcessingFinished int `json:"ProcessingFinished"`
	TransactionInfo    struct {
		ComputerName  string    `json:"ComputerName"`
		DateTime      time.Time `json:"DateTime"`
		SystemInfo    string    `json:"SystemInfo"`
		TransactionID string    `json:"TransactionID"`
		UserName      string    `json:"UserName"`
		Version       string    `json:"Version"`
	} `json:"TransactionInfo"`
	ElapsedTime        int         `json:"elapsedTime"`
	MorePagesAvailable int         `json:"morePagesAvailable"`
	PassBackObject     interface{} `json:"passBackObject"`
}

const (
	Address        = "Address"
	Authority      = "Authority"
	DateOfBirth    = "Date of Birth"
	DateOfExpiry   = "Date of Expiry"
	DateOfIssue    = "Date of Issue"
	DocumentNumber = "Document Number"
	Name           = "Surname And Given Names"
	Nationality    = "Nationality"
	Sex            = "Sex"
)

type Document struct {
	Address        string `json:"address,omitempty" example:"住址"`
	Authority      string `json:"authority,omitempty" example:"签发机关,"` // 24
	DateOfBirth    string `json:"date_of_birth,omitempty" example:"出生日期"`
	DateOfExpiry   string `json:"date_of_expiry,omitempty" example:"到期日期"` // 3
	DateOfIssue    string `json:"date_of_issue,omitempty" example:"签发日期"`  // 2
	DocumentNumber string `json:"document_number,omitempty" example:"证件号"`
	Name           string `json:"name,omitempty" example:"姓名"`
	Nationality    string `json:"nationality,omitempty" example:"国籍"`
	Sex            string `json:"sex,omitempty" example:"性别"`
}

func Read(path string) string {
	content, err := ioutil.ReadFile(path)

	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(content)
}

func main() {
	var response Response

	var image = Read("test.jpg")
	var payload = NewPayload(image)

	err := gout.POST("http://localhost:8080/api/process").
		SetJSON(payload).
		BindJSON(&response).
		Debug(false).
		Do()

	if err != nil {
		panic(err)
	}

	var document Document

	for _, item := range response.ContainerList.List {
		if item.ResultType == 36 {
			for _, field := range item.Text.FieldList {
				switch field.FieldName {
				case Address:
					document.Address = field.Value
				case Authority:
					document.Authority = field.Value
				case DateOfBirth:
					document.DateOfBirth = field.Value
				case DateOfExpiry:
					document.DateOfExpiry = field.Value
				case DateOfIssue:
					document.DateOfIssue = field.Value
				case DocumentNumber:
					document.DocumentNumber = field.Value
				case Name:
					document.Name = field.Value
				case Nationality:
					document.Nationality = field.Value
				case Sex:
					document.Sex = field.Value
				}
			}
		}
	}

	fmt.Printf("%+v", document)
}
