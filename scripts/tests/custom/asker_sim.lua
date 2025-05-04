-- Predefined arrays
local swiftCodes = {"AAISALTRXXX", "BPKOPLPWXXX", "BREXPLPWXXX", "BREXPLPWMBK", "BSCHCLR10R5"}
local countryCodes = {"BG", "PL", "MT", "CL","LV","UY"}
-- Base URL parts
local basePath = "/v1/swift-codes"
local requestCounter = 0
local lastPostedSwiftCode = nil

-- Set up before each thread starts
math.randomseed(os.time())

-- Choose method and path per request
request = function()
    local choice = math.random(1, 3)
    local swiftCode = swiftCodes[math.random(#swiftCodes)]
    local countryCode = countryCodes[math.random(#countryCodes)]
    
    -- If we posted a swift code last request, we need to delete it
    if lastPostedSwiftCode then
        local deleteRequest = wrk.format("DELETE", basePath .. "/" .. lastPostedSwiftCode)
        lastPostedSwiftCode = nil
        return deleteRequest
    end

    if choice == 1 then
        -- GET /swift-codes/{swiftCode}
        return wrk.format("GET", basePath .. "/" .. swiftCode)

    elseif choice == 2 then
        -- GET /swift-codes/country/{countryCode}
        return wrk.format("GET", basePath .. "/country/" .. countryCode)

    elseif choice == 3 then
        -- POST /swift-codes
        local newSwiftCode = "TESTDE" .. math.random(10, 99) .. "XXX"
        local body = [[
        {
            "countryISO2": "DE",
            "swiftCode": "]] .. newSwiftCode .. [[",
            "bankName": "THE GERMAN BANK",
            "address": "456 TEST STRASSE",
            "isHeadquarter": true,
            "countryName": "GERMANY"
        }]]
        
        -- Store the swift code to delete in the next request
        lastPostedSwiftCode = newSwiftCode
        return wrk.format("POST", basePath, {["Content-Type"] = "application/json"}, body)
    end
end

-- This function is called when a response is received
response = function(status, headers, body)
end