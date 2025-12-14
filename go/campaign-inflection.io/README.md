

Granularity
- second
- minute
- hour

UserProfile
- UserID string
- Email string
- PhoneNumber string

Campaign
- ID string
- Mode ENUM - Email / SMS / etc.
- List<UserProfile>
- IsRecurring bool
- RunExpression string "* * Tuesday 5 25"
- NextRunTimestamp Timestamp

CampaignExecutions
- CampaignID string FK
- ExecutionID UUID string
- status ENUM <InProgress, Success, Failed>

CampaignExecutionsUserProfile
- CampaignExecutionID string
- UserID string
- CampaignStatus -> Success / Failed <<- synchronous call response
- UserAction -> Clicked / Unsubscribed / Opened
- Meta JSON:
  actionHistory: [
    {
      "timestamp": 19190100,
      "action": "EmailClicked"
    }
  ]

ReportingData
- ID
- CampaignExecutionID
- UserID
- CampaignType (SMS/Email)
- EmailClicked bool
- EmailUnsubscribed bool
- ....



AnalyticsRepository
- Map<CampaignExecutionID, Map<UserProfile, ReportingData>> analyticsStore
- CRUD

CampaignRepository
- Map<ID, Campaign> campaignRepo // datastore
- sync.RWMutex
- CRUD()
- UpdateNextRunTimestamp(CampaignID string, NextRunTimestamp Timestamp) bool

Campaign Manager
- CampaignRepository
- CRUD()

- CRONs - second level granularity
- minute level granularity --> 1 worker

- GetAllCampaignsToRun(granularity Granularity, currentTs Timestamp) List<Campaign>, error
  - CampaignRepository -> tillNextTs: {currentTs + fetch granularity based window
  } -> fetch all campaigns where tillNextTs >= NextRunTimestamp
  - UpdateNextRunTimestamp is called for all returned Campaign

- RunAllCompaigns(List<Campaign>)
- RunCampaign(Campaign)
  - Add CampaignExecutions entry with InProgress
  - basis Mode, we can have strategy to execute Campaigns
    - Response:
      {
        "responseActions": [ { "userId": "", actionId: "", status: "Sent/NotSent", "failure_reason": ""} ]
      }
  - Loop through responseActions, check if all sent status i
  - save to AnalyticsRepository for list i.e. responseActions
    - campaignID 

- HandleUserAction(UserID, EmailID, CampaignExecutionID string) 
  - update CampaignExecutionsUserProfile with current action
  - in-memory ReportingData action to be updated to AnalyticsRepository
  - Get on UserProfile for UserID
  - Get&Update for CampaignExecutionID x UserProfile to update ReportingData in analyticsStore

Campaign Scheduler
