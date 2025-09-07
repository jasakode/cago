package cago

// Timezone is a placeholder type for potential future features that may
// depend on time zones. It is currently unused by the cache engine.
type Timezone string

const (
    // Universal
    TimezoneUTC   Timezone = "UTC"   // UTC+00
    TimezoneLocal Timezone = "Local" // host system local timezone

    // Indonesia
    TimezoneWIB  Timezone = "Asia/Jakarta"  // UTC+07
    TimezoneWITA Timezone = "Asia/Makassar" // UTC+08
    TimezoneWIT  Timezone = "Asia/Jayapura" // UTC+09

    // Asia
    TimezoneTokyo     Timezone = "Asia/Tokyo"     // UTC+09
    TimezoneSingapore Timezone = "Asia/Singapore" // UTC+08
    TimezoneBangkok   Timezone = "Asia/Bangkok"   // UTC+07
    TimezoneShanghai  Timezone = "Asia/Shanghai"  // UTC+08
    TimezoneDubai     Timezone = "Asia/Dubai"     // UTC+04

    // Europe
    TimezoneLondon Timezone = "Europe/London" // UTC+00 (DST +01)
    TimezoneBerlin Timezone = "Europe/Berlin" // UTC+01 (DST +02)
    TimezoneParis  Timezone = "Europe/Paris"  // UTC+01 (DST +02)
    TimezoneMoscow Timezone = "Europe/Moscow" // UTC+03

    // America
    TimezoneNewYork    Timezone = "America/New_York"    // UTC-05 (DST -04)
    TimezoneLosAngeles Timezone = "America/Los_Angeles" // UTC-08 (DST -07)
    TimezoneChicago    Timezone = "America/Chicago"     // UTC-06 (DST -05)
    TimezoneSaoPaulo   Timezone = "America/Sao_Paulo"   // UTC-03
    TimezoneMexicoCity Timezone = "America/Mexico_City" // UTC-06 (DST -05)

    // Australia / Pacific
    TimezoneSydney   Timezone = "Australia/Sydney" // UTC+10 (DST +11)
    TimezoneAuckland Timezone = "Pacific/Auckland" // UTC+12 (DST +13)
    TimezoneHonolulu Timezone = "Pacific/Honolulu" // UTC-10
)
