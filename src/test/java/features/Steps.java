package features;

import cucumber.api.java.en.And;
import cucumber.api.java.en.Given;
import cucumber.api.java.en.When;
import org.gatblau.onix.Info;
import org.gatblau.onix.Result;
import org.gatblau.onix.data.ItemData;
import org.gatblau.onix.data.ItemList;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpMethod;
import org.springframework.http.ResponseEntity;

import javax.annotation.PostConstruct;
import java.time.ZoneId;
import java.time.ZonedDateTime;
import java.time.format.DateTimeFormatter;
import java.util.HashMap;
import java.util.Map;

import static features.Key.*;

public class Steps extends BaseTest {
    private String baseUrl;

    @PostConstruct
    public void init() {
        baseUrl = "http://localhost:" + port + "/";
    }

    @And("^the base URL of the service is known$")
    public void theBaseURLOfTheServiceIsKnown() throws Throwable {
        util.put(BASE_URL, baseUrl);
    }

    @And("^a get request to the service is done$")
    public void aGetRequestToTheServiceIsDone() throws Throwable {
        ResponseEntity<Info> response = client.getForEntity((String)util.get(BASE_URL), Info.class);
        util.put(RESPONSE, response);
    }

    @And("^the service responds with description and version number$")
    public void theServiceRespondsWithDescriptionAndVersionNumber() throws Throwable {
        ResponseEntity<Info> response = util.get(RESPONSE);
        assert (response.getStatusCode().value() == 200);
    }

    @Given("^the item URL search by key is known$")
    public void theItemURLSearchByKeyIsKnown() throws Throwable {
        util.put(Key.ITEM_URL, String.format("%sitem/{key}/", baseUrl));
    }

    @Given("^the item URL search with query parameters is known$")
    public void theItemURLSearchWithQueryParametersIsKnown() throws Throwable {
        util.put(Key.ITEM_URL, String.format("%sitem/search", baseUrl));
    }

    @And("^the response code is (\\d+)$")
    public void theResponseCodeIs(int responseCode) throws Throwable {
        ResponseEntity<Result> response = util.get(RESPONSE);
        assert (response.getStatusCodeValue() == responseCode);
    }

    @And("^the response has body$")
    public void theResponseHasBody() throws Throwable {
        ResponseEntity<Result> response = util.get(RESPONSE);
        assert (response.hasBody());
    }

    @And("^a json payload with new item information exists$")
    public void aJsonPayloadWithNewItemInformationExists() throws Throwable {
        String payload = util.getFile("payload/create_item_payload.json");
        util.put(Key.PAYLOAD, payload);
    }

    @And("^the item does not exist in the database$")
    public void theItemDoesNotExistInTheDatabase() throws Throwable {
        theClearCMDBURLOfTheServiceIsKnown();
        aClearCMDBRequestToTheServiceIsDone();
    }

    @And("^the database is cleared$")
    public void theDatabaseIsCleared() throws Throwable {
        assert(!util.containsKey(EXCEPTION));
    }

    @And("^the clear cmdb URL of the service is known$")
    public void theClearCMDBURLOfTheServiceIsKnown() throws Throwable {
        util.put(Key.CLEAR_URL, String.format("%s/clear/", baseUrl));
    }

    @And("^a clear cmdb request to the service is done$")
    public void aClearCMDBRequestToTheServiceIsDone() throws Throwable {
        try {
            client.delete((String) util.get(CLEAR_URL));
            util.remove(EXCEPTION);
        }
        catch (Exception ex) {
            util.put(EXCEPTION, ex);
        }
    }

    @And("^there is not any error in the response$")
    public void thereIsNotAnyErrorInTheResponse() throws Throwable {
        assert(!util.containsKey(EXCEPTION));
    }

    @And("^a PUT HTTP request with a JSON payload is done$")
    public void aPUTTHTTPRequestWithAJSONPayloadIsDone() throws Throwable {
        putItem(ITEM_ONE_KEY, "payload/create_item_payload.json");
    }

    @And("^a configuration item natural key is known$")
    public void aConfigurationItemNaturalKeyIsKnown() throws Throwable {
        util.put(ITEM_ONE_KEY, "ITEM_ONE_KEY");
    }

    @And("^the service responds with action \"([^\"]*)\"$")
    public void theServiceRespondsWithAction(String action) throws Throwable {
        ResponseEntity<Result> response = util.get(RESPONSE);
        assert (response.getBody().getAction().equals(action));
    }

    @And("^the item exist in the database$")
    public void theItemExistInTheDatabase() throws Throwable {
        theItemDoesNotExistInTheDatabase();
        theItemURLSearchByKeyIsKnown();
        aJsonPayloadWithNewItemInformationExists();
        aPUTTHTTPRequestWithAJSONPayloadIsDone();
    }

    @And("^a json payload with updated item information exists$")
    public void aJsonPayloadWithUpdatedItemInformationExists() throws Throwable {
        String payload = util.getFile("payload/update_item_payload.json");
        util.put(Key.PAYLOAD, payload);
    }

    @Given("^the item type does not exist in the database$")
    public void theItemTypeDoesNotExistInTheDatabase() throws Throwable {
        theItemTypeURLOfTheServiceIsKnown();
        aDELETEHTTPRequestIsDone();
        thereIsNotAnyErrorInTheResponse();
    }

    @Given("^the item type URL of the service is known$")
    public void theItemTypeURLOfTheServiceIsKnown() throws Throwable {
        util.put(ITEM_TYPE_URL, String.format("%s/itemtype/", baseUrl));
    }

    @When("^a DELETE HTTP request is done$")
    public void aDELETEHTTPRequestIsDone() throws Throwable {
        try {
            client.delete((String) util.get(ITEM_TYPE_URL));
            util.remove(EXCEPTION);
        }
        catch (Exception ex) {
            util.put(EXCEPTION, ex);
        }
    }

    @Given("^the link URL of the service is known$")
    public void theLinkURLOfTheServiceIsKnown() throws Throwable {
        util.put(LINK_URL, String.format("%s/link/{key}/", baseUrl));
    }

    @Given("^a json payload with new link information exists$")
    public void aJsonPayloadWithNewLinkInformationExists() throws Throwable {
        String payload = util.getFile("payload/create_link_payload.json");
        util.put(Key.PAYLOAD, payload);
    }

    @Given("^a link to the two configuration items does not exist in the database$")
    public void aLinkToTheTwoConfigurationItemsDoesNotExistInTheDatabase() throws Throwable {
        aDELETELinkRequestIsDone();
    }

    @When("^a DELETE Link request is done$")
    public void aDELETELinkRequestIsDone() throws Throwable {
        try {
            client.delete((String) util.get(LINK_URL), (String)util.get(LINK_KEY));
            util.remove(EXCEPTION);
        }
        catch (Exception ex) {
            util.put(EXCEPTION, ex);
        }
    }

    @When("^a PUT HTTP request to the Link resource is done with a JSON payload$")
    public void aPUTHTTPRequestToTheLinkResourceIsDoneWithAJSONPayload() throws Throwable {
        String url = util.get(LINK_URL);
        Map<String, Object> vars = new HashMap<>();
        vars.put("key", util.get(LINK_KEY));
        ResponseEntity<Result> response = null;
        try {
            response = client.exchange(url, HttpMethod.PUT, getEntity(PAYLOAD), Result.class, vars);
            util.put(RESPONSE, response);
            util.remove(EXCEPTION);
        }
        catch (Exception ex) {
            util.put(EXCEPTION, ex);
        }
    }

    private HttpEntity<?> getEntity(String payloadKey) {
        String payload = util.get(payloadKey);
        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
       return new HttpEntity<>(payload, headers);
    }

    @Given("^the configuration items to be linked exist in the database$")
    public void theConfigurationItemsToBeLinkedExistInTheDatabase() throws Throwable {
        putItem(ITEM_ONE_KEY, "payload/create_item_payload.json");
        putItem(ITEM_TWO_KEY, "payload/create_item_payload.json");
    }

    @Given("^the item exists in the database$")
    public void theItemExistsInTheDatabase() throws Throwable {
        putItem("item_one", "payload/update_item_payload.json");
    }

    @When("^a GET HTTP request to the Item uri is done$")
    public void aGETHTTPRequestToTheItemUriIsDone() throws Throwable {
        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
        ResponseEntity<ItemData> result = client.exchange(
                (String)util.get(ITEM_URL),
                HttpMethod.GET,
                new HttpEntity<>(null, headers),
                ItemData.class,
                (String)util.get(ITEM_ONE_KEY));
        util.put(RESPONSE, result);
    }

    private void putItem(String itemKey, String filename) {
        util.put(Key.PAYLOAD, util.getFile(filename));
        String url = String.format("%s/item/{key}/", baseUrl);
        Map<String, Object> vars = new HashMap<>();
        vars.put("key", itemKey);
        ResponseEntity<Result> response = null;
        try {
            response = client.exchange(url, HttpMethod.PUT, getEntity(PAYLOAD), Result.class, vars);
            util.put(RESPONSE, response);
            util.remove(EXCEPTION);
        }
        catch (Exception ex) {
            util.put(EXCEPTION, ex);
        }
    }

    @Given("^more than one item exist in the database$")
    public void moreThanOneItemExistInTheDatabase() throws Throwable {
        putItem("item_one", "payload/update_item_payload.json");
        putItem("item_two", "payload/update_item_payload.json");
        putItem("item_three", "payload/update_item_payload.json");
    }

    @When("^a GET HTTP request to the Item uri is done with query parameters$")
    public void aGETHTTPRequestToTheItemUriIsDoneWithQueryParameters() throws Throwable {
        DateTimeFormatter formatter = DateTimeFormatter.ofPattern("yyyyMMddHHmm");

        StringBuilder uri = new StringBuilder();
        uri.append((String)util.get(ITEM_URL));

        if (util.containsKey(CONGIG_ITEM_TYPE_ID) || util.containsKey(CONFIG_ITEM_TAG) || util.containsKey(CONFIG_ITEM_UPDATED_FROM)) {
            uri.append("?");
        }

        if (util.containsKey(CONGIG_ITEM_TYPE_ID)) {
            Integer typeId = util.get(CONGIG_ITEM_TYPE_ID);
            uri.append("typeId=").append(typeId).append("&");
        }

        if (util.containsKey(CONFIG_ITEM_TAG)) {
            String tag = util.get(CONFIG_ITEM_TAG);
            uri.append("tag=").append(tag).append("&");
        }

        if (util.containsKey(CONFIG_ITEM_UPDATED_FROM)) {
            ZonedDateTime from = util.get(CONFIG_ITEM_UPDATED_FROM);
            uri.append("from=").append(from.format(formatter)).append("&");
        }

        if (util.containsKey(CONFIG_ITEM_UPDATED_TO)) {
            ZonedDateTime to = util.get(CONFIG_ITEM_UPDATED_TO);
            uri.append("to=").append(to.format(formatter)).append("&");
        }

        String uriString = uri.toString();
        if (uriString.endsWith("&")) {
            uriString = uriString.substring(0, uriString.length() - 1);
        }

        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
        ResponseEntity<ItemList> result = client.exchange(
                uriString,
                HttpMethod.GET,
                new HttpEntity<>(null, headers),
                ItemList.class);
        util.put(RESPONSE, result);
    }

    @Given("^the filtering config item type is known$")
    public void theFilteringConfigItemTypeIsKnown() throws Throwable {
        util.put(CONGIG_ITEM_TYPE_ID, 2);
    }

    @Given("^the filtering config item tag is known$")
    public void theFilteringConfigItemTagIsKnown() throws Throwable {
        util.put(CONFIG_ITEM_TAG, "Test");
    }

    @Given("^the filtering config item date range is known$")
    public void theFilteringConfigItemDateRangeIsKnown() throws Throwable {
        util.put(CONFIG_ITEM_UPDATED_FROM, ZonedDateTime.of(ZonedDateTime.now().getYear() - 100, 1, 1, 0, 0, 0, 0, ZoneId.systemDefault()));
        util.put(CONFIG_ITEM_UPDATED_TO, ZonedDateTime.of(ZonedDateTime.now().getYear(), 1, 1, 0, 0, 0, 0, ZoneId.systemDefault()));
    }

    @Given("^the natural key for the link is known$")
    public void theNaturalKeyForTheLinkIsKnown() throws Throwable {
        util.put(LINK_KEY, "LINK_KEY");
    }
}
