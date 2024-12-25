package moneyforward

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type MFPath string // eg "sp2/accounts/t0qRlCziUbsxYAgcH2fGbw/edit"

type MFShowPath MFPath

// HomeTimelineResponse represents the response from the home timeline endpoint
type HomeTimelineResponse struct {
	Self struct {
		Href  string `json:"href"`
		Path  string `json:"path"`
		Until string `json:"until"`
		Limit int    `json:"limit"`
	} `json:"self"`
	RequestTime time.Time `json:"request_time"`
	Timeline    []struct {
		Date  string `json:"date"`
		Cards []struct {
			Type             string `json:"type"`
			UserNotification *struct {
				ID         int64 `json:"id"`
				CategoryID int   `json:"category_id"`
				Category   struct {
					PremiumRequired bool `json:"premium_required"`
				} `json:"category"`
				Parameters struct {
					Account *struct {
						Name                 string     `json:"name"`
						Amount               float64    `json:"amount"`
						Status               int        `json:"status"`
						Type                 string     `json:"type"`
						ShowPath             MFShowPath `json:"show_path"`
						LastSucceededAt      string     `json:"last_succeeded_at"`
						AggregationQueuePath MFPath     `json:"aggregation_queue_path"`
						AccountIDHash        string     `json:"account_id_hash"`
						ServiceID            int        `json:"service_id"`
						ServiceType          string     `json:"service_type"`
						ServiceCategoryID    int        `json:"service_category_id"`
						ColorCode            string     `json:"color_code"`
						IsShowTransaction    bool       `json:"is_show_transaction"`
					} `json:"account,omitempty"`
					UserAssetActIDs []int64 `json:"user_asset_act_ids"`
					LargestAmount   *struct {
						Amount float64   `json:"amount"`
						Date   time.Time `json:"date"`
					} `json:"largest_amount"`
					SumAmount *struct {
						Amount float64   `json:"amount"`
						Date   time.Time `json:"date"`
					} `json:"sum_amount"`
					Extra map[string]interface{} `json:"extra"`
				} `json:"parameters"`
				ReadAt *time.Time `json:"read_at"`
				Read   bool       `json:"read"`
			} `json:"user_notification,omitempty"`
			HomeCard *struct {
				ID          string `json:"id"`
				BannerImage struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"banner_image"`
				LandingURL string    `json:"landing_url"`
				CreatedAt  time.Time `json:"created_at"`
				StartAt    time.Time `json:"start_at"`
				EndAt      time.Time `json:"end_at"`
			} `json:"home_card,omitempty"`
		} `json:"cards"`
	} `json:"timeline"`
}

// AccountSummariesResponse represents the response from the account summaries endpoint
type AccountSummariesResponse struct {
	Accounts []struct {
		Name                 string     `json:"name"`
		Amount               float64    `json:"amount"`
		LastLoginAt          string     `json:"last_login_at"`
		LastAggregatedAt     string     `json:"last_aggregated_at"`
		LastSucceededAt      string     `json:"last_succeeded_at"`
		ErrorID              int        `json:"error_id"`
		Status               int        `json:"status"`
		Type                 string     `json:"type"`
		AccountIDHash        string     `json:"account_id_hash"`
		ShowPath             MFShowPath `json:"show_path"`
		AggregationQueuePath MFPath     `json:"aggregation_queue_path"`
		ServiceID            int        `json:"service_id"`
		ServiceType          string     `json:"service_type"`
		ServiceCategoryID    int        `json:"service_category_id"`
		SubAccounts          []struct {
			SubAccountIDHash      string `json:"sub_account_id_hash"`
			SubName               string `json:"sub_name"`
			SubType               string `json:"sub_type"`
			SubNumber             string `json:"sub_number"`
			UserAssetDetSummaries []struct {
				AssetClassID      int     `json:"asset_class_id"`
				AssetSubclassID   int     `json:"asset_subclass_id"`
				AssetSubclassName string  `json:"asset_subclass_name"`
				AssetSubclassUnit string  `json:"asset_subclass_unit"`
				Value             float64 `json:"value"`
				JPYValue          float64 `json:"jpyvalue"`
			} `json:"user_asset_det_summaries"`
		} `json:"sub_accounts"`
	} `json:"accounts"`
}

// TransactionsResponse represents the response from the transactions endpoint
type TransactionsResponse struct {
	EmptyState struct {
		RecommendedServices struct {
			Services []struct {
				ID              int    `json:"id"`
				ServiceName     string `json:"service_name"`
				ServiceType     string `json:"service_type"`
				ColorCode       string `json:"color_code"`
				ServiceCategory struct {
					ID           int    `json:"id"`
					CategoryType string `json:"category_type"`
				} `json:"service_category"`
			} `json:"services"`
		} `json:"recommended_services"`
	} `json:"empty_state"`
}

// UserAssetActsResponse represents the response from the user asset acts endpoint
type UserAssetActsResponse struct {
	UserAssetActs  []*UserAssetAct `json:"user_asset_acts"`
	RecordCount    int             `json:"record_count"`
	TotalCount     int             `json:"total_count"`
	Offset         int             `json:"offset"`
	Size           int             `json:"size"`
	From           string          `json:"from"`
	To             string          `json:"to"`
	NewRecordCount int             `json:"new_record_count"`
}

// UserAssetActResponse represents the response for a single user asset act
type UserAssetActResponse struct {
	UserAssetAct UserAssetAct `json:"user_asset_act"`
}

// StringID is a custom type that can unmarshal both string and number JSON values into a string
type StringID string

func (sid *StringID) String() string {
	return string(*sid)
}

func (sid *StringID) UnmarshalJSON(data []byte) error {
	// Try string first
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*sid = StringID(s)
		return nil
	}

	// Try number
	var i int64
	if err := json.Unmarshal(data, &i); err == nil {
		*sid = StringID(strconv.FormatInt(i, 10))
		return nil
	}

	return fmt.Errorf("value must be string or number, got %s", data)
}

// UserAssetAct represents a single user asset activity
type UserAssetAct struct {
	ID               StringID  `json:"id"`
	AccountID        StringID  `json:"account_id"`
	SubAccountID     StringID  `json:"sub_account_id"`
	IsTransfer       bool      `json:"is_transfer"`
	IsIncome         bool      `json:"is_income"`
	Content          string    `json:"content"`
	OrigContent      string    `json:"orig_content"`
	Amount           float64   `json:"amount"`
	OrigAmount       float64   `json:"orig_amount"`
	Currency         string    `json:"currency"`
	JPYRate          float64   `json:"jpyrate"`
	LargeCategoryID  StringID  `json:"large_category_id"`
	MiddleCategoryID StringID  `json:"middle_category_id"`
	CreatedAt        time.Time `json:"created_at"`
	RecognizedAt     time.Time `json:"recognized_at"`
	UpdatedAt        string    `json:"updated_at"`
	Account          struct {
		ServiceID         StringID `json:"service_id"`
		ServiceCategoryID StringID `json:"service_category_id"`
		Service           struct {
			ServiceName string `json:"service_name"`
		} `json:"service"`
	} `json:"account"`
	SubAccount struct {
		SubName   string `json:"sub_name"`
		SubType   string `json:"sub_type"`
		SubNumber string `json:"sub_number"`
	} `json:"sub_account"`
	IsJournalizable bool `json:"is_journalizable_service"`
	IsJournalized   bool `json:"is_journalized"`
}

// AccountResponse represents the response from the account endpoint
type AccountResponse struct {
	Account struct {
		ServiceID         int       `json:"service_id"`
		Status            string    `json:"status"`
		ErrorID           string    `json:"error_id"`
		LastLoginAt       time.Time `json:"last_login_at"`
		LastSucceededAt   string    `json:"last_succeeded_at"`
		LastAggregatedAt  time.Time `json:"last_aggregated_at"`
		AccountIDHash     string    `json:"account_id_hash"`
		DisplayName       string    `json:"display_name"`
		ServiceCategoryID string    `json:"service_category_id"`
		TotalAsset        float64   `json:"total_asset"`
		TotalLiability    float64   `json:"total_liability"`
		SubAccounts       []struct {
			SubAccountIDHash      string `json:"sub_account_id_hash"`
			SubName               string `json:"sub_name"`
			SubType               string `json:"sub_type"`
			SubNumber             string `json:"sub_number"`
			ServiceCategoryID     string `json:"service_category_id"`
			UserAssetDetSummaries []struct {
				AssetClassID      int     `json:"asset_class_id"`
				AssetSubclassID   int     `json:"asset_subclass_id"`
				AssetSubclassName string  `json:"asset_subclass_name"`
				Value             float64 `json:"value"`
				JPYValue          float64 `json:"jpyvalue"`
			} `json:"user_asset_det_summaries"`
		} `json:"sub_accounts"`
		Service struct {
			ServiceType           string `json:"service_type"`
			LoginURL              string `json:"login_url"`
			IsShowTransaction     bool   `json:"is_show_transaction"`
			ColorCode             string `json:"color_code"`
			Aggregable            bool   `json:"aggregable"`
			RequiresUserOperation bool   `json:"requires_user_operation"`
		} `json:"service"`
	} `json:"account"`
}

// AssetHistoryResponse represents the asset history response
type AssetHistoryResponse struct {
	Histories []struct {
		Date      string  `json:"date"`
		Amount    float64 `json:"amount"`
		Category  string  `json:"category"`
		ClassID   int     `json:"class_id"`
		ClassName string  `json:"class_name"`
	} `json:"histories"`
}

// AssetClassesResponse represents the asset classes response
type AssetClassesResponse struct {
	AssetClasses []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"asset_classes"`
}

// AssetSubclassesResponse represents the asset subclasses response
type AssetSubclassesResponse struct {
	AssetSubclasses []struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		ClassID  int    `json:"class_id"`
		Unit     string `json:"unit"`
		Ordering int    `json:"ordering"`
	} `json:"asset_subclasses"`
}

// CategoriesResponse represents the categories response
type CategoriesResponse struct {
	LargeCategories []struct {
		ID               int    `json:"id"`
		Name             string `json:"name"`
		MiddleCategories []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"middle_categories"`
	} `json:"large_categories"`
}

// UserAssetActUpdates represents the update payload for transactions
type UserAssetActUpdates struct {
	LargeCategoryID  int     `json:"large_category_id,omitempty"`
	MiddleCategoryID int     `json:"middle_category_id,omitempty"`
	Content          string  `json:"content,omitempty"`
	Amount           float64 `json:"amount,omitempty"`
}

// ServiceCategoriesResponse represents the service categories response
type ServiceCategoriesResponse struct {
	Categories []struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		ServiceType string `json:"service_type"`
	} `json:"categories"`
}

// ServicesResponse represents the services response
type ServicesResponse struct {
	Services []struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		ServiceType string `json:"service_type"`
		CategoryID  int    `json:"category_id"`
	} `json:"services"`
}

// CashFlowTermDataResponse represents the response from the cash flow term data endpoint
type CashFlowTermDataResponse struct {
	Result        string `json:"result"`
	UserAssetActs []struct {
		UserAssetAct UserAssetAct `json:"user_asset_act"`
	} `json:"user_asset_acts"`
}
