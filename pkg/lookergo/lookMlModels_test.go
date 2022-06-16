package lookergo

/*

### Get All LookMlModels
GET {{endpoint}}/4.0/lookml_models
Authorization: token {{token}}
`[
  {
    "has_content": false,
    "label": "Reprise J&J",
    "name": "Reprise_J&J",
    "project_name": "us_johnson_johnson",
    "unlimited_db_connections": false,
    "allowed_db_connection_names": [
      "data-mesh-googleanalytics-block"
    ],
    "explores": [],
    "can": {
      "index": true,
      "show": true,
      "create": true,
      "update": true,
      "destroy": true
    }
  },
  {
    "has_content": true,
    "label": "Cross Channel",
    "name": "cross_channel",
    "project_name": "hub",
    "unlimited_db_connections": true,
    "allowed_db_connection_names": [
      "trial_bigquery_aussiehomeloans",
      "data-mesh-search_ads_360-block",
      "data-mesh-googleanalytics-block",
      "data-mesh-google-ads",
      "data-mesh-google-ads-de",
      "data-mesh-search_ads_360-block_de",
      "sven-test",
      "de_facebookads_aida",
      "data-mesh-shark-ninja",
      "johnson-johnson-connection-test3",
      "looker_test_reprise_th",
      "data-mesh-google-ads-officeworks",
      "bq-looker-connection-dach-cross-client",
      "bq-looker-connection-global-prime-video",
      "us_ecomm_reviews",
      "pj-reprise-billing-export",
      "bq-looker-connection-usa-sharkninja",
      "bq-looker-connection-usa-levistrauss",
      "bq-looker-connection-usa-j3nutragena",
      "bq-looker-connection-usa-aldi",
      "bq-looker-connection-usa-aldius",
      "bq-looker-connection-nld-hbm",
      "bq-looker-connection-dach-amazon-wfs",
      "bq-looker-connection-sandbox-jonas-fresh",
      "bq-looker-connection-sandbox-vincent1",
      "bq-looker-connection-sandbox-vincent2",
      "bq-looker-connection-sandbox-frank-1",
      "bq-looker-connection-learning-customer-1",
      "bq-looker-connection-usa-leafguard",
      "bq-looker-connection-sandbox-jonas",
      "bq-looker-connection-dach-boconcept-de",
      "bq-looker-connection-dach-boconcept",
      "bq-looker-connection-dach-dyson-ch",
      "bq-looker-connection-dach-dyson-de",
      "bq-looker-connection-canada-amz-xcm",
      "bq-looker-connection-sandbox-emea-jonas",
      "bq-looker-connection-uk-tommee-t",
      "bq-looker-connection-dach-teva-de",
      "bq-looker-connection-usa-amazon-gca",
      "bq-looker-connection-dach-ch",
      "bq-looker-connection-nld-plussuper",
      "bq-looker-connection-dach-playground",
      "bq-looker-connection-sandbox-bertrand2",
      "bq-looker-connection-usa-x-client",
      "bq-looker-connection-mexico-unilever",
      "bq-looker-connection-canada-amz-prime-video",
      "bq-looker-connection-canada-sandbox",
      "bq-looker-connection-sandbox-flow",
      "bq-looker-connection-usa-fda-vaping",
      "bq-looker-connection-usa-sandbox",
      "bq-looker-connection-usa-amz-primevideo",
      "bq-looker-connection-usa-quickbooks",
      "bq-looker-connection-uk-dyson",
      "bq-looker-connection-nld-amex",
      "bq-looker-connection-dach-vzug-ch"
    ],
    "explores": [
      {
        "description": "Cross-channel explore view that supports a sub-set of the channel specific measures and dimensions.",
        "label": "Cross Channel ❎ ",
        "hidden": false,
        "group_label": "hub Social Channels",
        "name": "cross_channel"
      }
    ],
    "can": {
      "index": true,
      "show": true,
      "create": true,
      "update": true,
      "destroy": true
    }
  },
  {
    "has_content": true,
    "label": "Amazon Prime",
    "name": "amazon_prime",
    "project_name": "amazon_prime",
    "unlimited_db_connections": true,
    "allowed_db_connection_names": [
      "trial_bigquery_aussiehomeloans",
      "data-mesh-search_ads_360-block",
      "data-mesh-googleanalytics-block",
      "data-mesh-google-ads",
      "data-mesh-google-ads-de",
      "data-mesh-search_ads_360-block_de",
      "sven-test",
      "de_facebookads_aida",
      "data-mesh-shark-ninja",
      "johnson-johnson-connection-test3",
      "looker_test_reprise_th",
      "data-mesh-google-ads-officeworks",
      "bq-looker-connection-dach-cross-client",
      "bq-looker-connection-global-prime-video",
      "us_ecomm_reviews",
      "pj-reprise-billing-export",
      "bq-looker-connection-usa-sharkninja",
      "bq-looker-connection-usa-levistrauss",
      "bq-looker-connection-usa-j3nutragena",
      "bq-looker-connection-usa-aldi",
      "bq-looker-connection-usa-aldius",
      "bq-looker-connection-nld-hbm",
      "bq-looker-connection-dach-amazon-wfs",
      "bq-looker-connection-sandbox-jonas-fresh",
      "bq-looker-connection-sandbox-vincent1",
      "bq-looker-connection-sandbox-vincent2",
      "bq-looker-connection-sandbox-frank-1",
      "bq-looker-connection-learning-customer-1",
      "bq-looker-connection-usa-leafguard",
      "bq-looker-connection-sandbox-jonas",
      "bq-looker-connection-dach-boconcept-de",
      "bq-looker-connection-dach-boconcept",
      "bq-looker-connection-dach-dyson-ch",
      "bq-looker-connection-dach-dyson-de",
      "bq-looker-connection-canada-amz-xcm",
      "bq-looker-connection-sandbox-emea-jonas",
      "bq-looker-connection-uk-tommee-t",
      "bq-looker-connection-dach-teva-de",
      "bq-looker-connection-usa-amazon-gca",
      "bq-looker-connection-dach-ch",
      "bq-looker-connection-nld-plussuper",
      "bq-looker-connection-dach-playground",
      "bq-looker-connection-sandbox-bertrand2",
      "bq-looker-connection-usa-x-client",
      "bq-looker-connection-mexico-unilever",
      "bq-looker-connection-canada-amz-prime-video",
      "bq-looker-connection-canada-sandbox",
      "bq-looker-connection-sandbox-flow",
      "bq-looker-connection-usa-fda-vaping",
      "bq-looker-connection-usa-sandbox",
      "bq-looker-connection-usa-amz-primevideo",
      "bq-looker-connection-usa-quickbooks",
      "bq-looker-connection-uk-dyson",
      "bq-looker-connection-nld-amex",
      "bq-looker-connection-dach-vzug-ch"
    ],
    "explores": [
      {
        "description": "Explores Facebook Ads data, broken down by the type of device, mobile or desktop, used by people when they viewed or clicked on an ad, as shown in ads reporting.",
        "label": "Facebook Ads: Device breakdown 👤",
        "hidden": false,
        "group_label": "amazon_prime Social Channels",
        "name": "facebook_ads_breakdown_device"
      },
      {
        "description": "Explores Facebook Ads data, broken down by age and gende of people you've reached. People who don't list their gender are shown as 'not specified'.",
        "label": "Facebook Ads: Age/Gender breakdown 👤",
        "hidden": false,
        "group_label": "amazon_prime Social Channels",
        "name": "facebook_ads_breakdown_age_gender"
      },
      {
        "description": "Explores Facebook data without breakdowns. Use this if you want to look at budget, ROAS or pacing.",
        "label": "Facebook Ads (no breakdown) 👤",
        "hidden": false,
        "group_label": "amazon_prime Social Channels",
        "name": "facebook_ads"
      },
      {
        "description": "Explores Snapchat data.",
        "label": "Snapchat 👻",
        "hidden": false,
        "group_label": "amazon_prime Social Channels",
        "name": "snapchat"
      },
      {
        "description": "Explores TikTok data.",
        "label": "TikTok 🎬",
        "hidden": false,
        "group_label": "amazon_prime Social Channels",
        "name": "tiktok"
      },
      {
        "description": null,
        "label": "Twitter Ads (No Breakdown) 🕊",
        "hidden": false,
        "group_label": "amazon_prime Social Channels",
        "name": "twitter_ads"
      },
      {
        "description": "Explores Twitter Ads data broken down by age.",
        "label": "Twitter Ads (Breakdown: Age) 🕊",
        "hidden": false,
        "group_label": "amazon_prime Social Channels",
        "name": "twitter_ads_breakdown_age"
      },
      {
        "description": "Explores Twitter Ads data broken down by gender.",
        "label": "Twitter Ads (Breakdown: Gender) 🕊",
        "hidden": false,
        "group_label": "amazon_prime Social Channels",
        "name": "twitter_ads_breakdown_gender"
      },
      {
        "description": "Explores Twitter Ads data broken down by device.",
        "label": "Twitter Ads (Breakdown: Device) 🕊",
        "hidden": false,
        "group_label": "amazon_prime Social Channels",
        "name": "twitter_ads_breakdown_device"
      },
      {
        "description": "Explores YouTube data coming through DV360.",
        "label": "YouTube (through DV360) 📺",
        "hidden": false,
        "group_label": "amazon_prime Social Channels",
        "name": "youtube_dv360"
      },
      {
        "description": "Explores YouTube data coming through Google Ads.",
        "label": "YouTube (through Google ads) 📺",
        "hidden": false,
        "group_label": "amazon_prime Social Channels",
        "name": "youtube_google_ads"
      },
      {
        "description": null,
        "label": "Cross Youtube Explore",
        "hidden": false,
        "group_label": "Amazon Prime",
        "name": "cross_youtube_sources"
      },
      {
        "description": null,
        "label": "Cross Youtube Explore (Breakdown: Age)",
        "hidden": false,
        "group_label": "Amazon Prime",
        "name": "cross_youtube_age"
      },
      {
        "description": null,
        "label": "Cross Youtube Explore (Breakdown: Gender)",
        "hidden": false,
        "group_label": "Amazon Prime",
        "name": "cross_youtube_gender"
      },
      {
        "description": "Cross-channel explore view that supports a sub-set of the channel specific measures and dimensions.",
        "label": "Cross Channel ❎ ",
        "hidden": false,
        "group_label": "amazon_prime Social Channels",
        "name": "cross_channel"
      }
    ],
    "can": {
      "index": true,
      "show": true,
      "create": true,
      "update": true,
      "destroy": true
    }
  },`

*/
