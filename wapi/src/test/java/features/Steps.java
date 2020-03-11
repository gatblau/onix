package features;

import cucumber.api.java.en.And;
import cucumber.api.java.en.Given;
import cucumber.api.java.en.Then;
import cucumber.api.java.en.When;
import org.gatblau.onix.data.*;
import org.json.simple.JSONObject;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpMethod;
import org.springframework.http.ResponseEntity;

import javax.annotation.PostConstruct;
import java.io.Serializable;
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

    @And("^a get request to the live url is done$")
    public void aGetRequestToTheServiceIsDone() throws Throwable {
        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "text/html");
        ResponseEntity<String> result = client.exchange(
                (String)util.get(LIVE_URL),
                HttpMethod.GET,
                new HttpEntity<>(null, headers),
                String.class);
        util.put(RESPONSE, result);
    }

    @Given("^the item URL search by key is known$")
    public void theItemURLSearchByKeyIsKnown() {
        util.put(ENDPOINT_URI, String.format("%sitem/{key}", baseUrl));
    }

    @Given("^the item URL search with query parameters is known$")
    public void theItemURLSearchWithQueryParametersIsKnown() {
        util.put(ITEM_URL, String.format("%s/item", baseUrl));
    }

    @And("^the response code is (\\d+)$")
    public void theResponseCodeIs(int responseCode)  {
        if (util.containsKey(EXCEPTION)) {
            RuntimeException ex = util.get(EXCEPTION);
            throw ex;
        }
        ResponseEntity<Result> response = util.get(RESPONSE);
        if (response.getStatusCode().value() != responseCode) {
            throw new RuntimeException(
                String.format(
                    "Expected response code was '%s' but instead got '%s': '%s'.",
                    responseCode,
                    response.getStatusCode(),
                    response.getStatusCode().getReasonPhrase()
                )
            );
        };
    }

    @And("^the response has body$")
    public void theResponseHasBody() throws Throwable {
        ResponseEntity<Result> response = util.get(RESPONSE);
        if (!response.hasBody()) {
            throw new RuntimeException("The response does not have a body.");
        };
    }

    @And("^a json payload with new item information exists$")
    public void aJsonPayloadWithNewItemInformationExists() throws Throwable {
        String payload = util.getFile("payload/create_item_payload.json");
        util.put(Key.PAYLOAD, payload);
    }

    @And("^the item does not exist in the database$")
    public void theItemDoesNotExistInTheDatabase() throws Throwable {
        Result result = delete(String.format("%s/item", baseUrl), null);
        if (result.isError()){
            throw new RuntimeException(result.getMessage());
        }
    }

    @And("^the database is cleared$")
    public void theDatabaseIsCleared() throws Throwable {
        if (util.containsKey(EXCEPTION)) {
            throw new RuntimeException((Exception)util.get(EXCEPTION));
        };
    }

    @And("^the clear cmdb URL of the service is known$")
    public void theClearCMDBURLOfTheServiceIsKnown() throws Throwable {
        util.put(CLEAR_URL, String.format("%s/clear", baseUrl));
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
    public void thereIsNotAnyErrorInTheResponse() {
        if (util.containsKey(RESPONSE)) {
            try {
                ResponseEntity<Result> response = util.get(RESPONSE);
                if (response.getBody() != null) {
                    if (response.getBody().isError()) {
                        throw new RuntimeException(response.getBody().getMessage());
                    }
                }
            } catch (ClassCastException cce) {
                System.out.println("WARNING: response type is not a Result. ");
            }
        }
        if (util.containsKey(EXCEPTION)){
            throw new RuntimeException((Exception)util.get(EXCEPTION));
        }
    }

    @And("^a PUT HTTP request with a new JSON payload is done$")
    public void aPUTTHTTPRequestWithANewJSONPayloadIsDone() throws Throwable {
        putItem(util.get(ITEM_ONE_KEY), "payload/create_item_payload.json");
    }

    @And("^a configuration item natural key is known$")
    public void aConfigurationItemNaturalKeyIsKnown() throws Throwable {
        util.put(ITEM_ONE_KEY, ITEM_ONE_KEY);
    }

    @And("^the service responds with action \"([^\"]*)\"$")
    public void theServiceRespondsWithAction(String action) throws Throwable {
        ResponseEntity<Result> response = util.get(RESPONSE);
        Result result = response.getBody();
        if (!result.getOperation().equals(action)) {
            throw new RuntimeException(String.format("Required result operation was %s but found %s", action, result.getOperation()));
        };
    }

    @And("^the item exist in the database$")
    public void theItemExistInTheDatabase() throws Throwable {
        theItemDoesNotExistInTheDatabase();
        theItemURLSearchByKeyIsKnown();
        aJsonPayloadWithNewItemInformationExists();
        aPUTTHTTPRequestWithANewJSONPayloadIsDone();
    }

    @Given("^the item type does not exist in the database$")
    public void theItemTypeDoesNotExistInTheDatabase() {
        delete(String.format("%sitemtype/{item_type}", baseUrl)+"?force=true", "item_type_1");
    }

    @Given("^the item type URL of the service is known$")
    public void theItemTypeURLOfTheServiceIsKnown() throws Throwable {
        util.put(ITEM_TYPE_URL, String.format("%sitemtype", baseUrl));
    }

    @Given("^the item type URL of the service with key is known$")
    public void theItemTypeURLOfTheServiceWithKeyIsKnown() throws Throwable {
        util.put(ITEM_TYPE_URL, String.format("%sitemtype/{key}", baseUrl));
    }

    @When("^a DELETE HTTP request with an item key is done$")
    public void aDELETEHTTPRequestWithAnItemKeyKeyIsDone() throws Throwable {
        Result result = delete(util.get(ENDPOINT_URI), util.get(ITEM_KEY));
        if (result.isError()){
            throw new RuntimeException(result.getMessage());
        }
    }

    @Given("^the link URL of the service is known$")
    public void theLinkURLOfTheServiceIsKnown() throws Throwable {
        util.put(LINK_URL, String.format("%s/link/{key}", baseUrl));
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
            response = client.exchange(url, HttpMethod.PUT, getEntityFromKey(PAYLOAD), Result.class, vars);
            util.put(RESPONSE, response);
            util.remove(EXCEPTION);
        }
        catch (Exception ex) {
            util.put(EXCEPTION, ex);
        }
    }

    @Given("^the configuration items to be linked exist in the database$")
    public void theConfigurationItemsToBeLinkedExistInTheDatabase() throws Throwable {
        putModel("meta_model_1", "payload/create_model_1_payload.json");
        putItemType("item_type_1", "payload/create_item_type_1_payload.json");
        putItemTypeAttr("item_type_1", "item_type_1_COMPANY", "payload/create_item_type_attr_1_payload.json");
        putItemTypeAttr("item_type_1", "item_type_1_WBS", "payload/create_item_type_attr_2_payload.json");
        putItem(ITEM_ONE_KEY, "payload/create_item_2_payload.json");
        putItem(ITEM_TWO_KEY, "payload/create_item_payload.json");
    }

    @Given("^the item exists in the database$")
    public void theItemExistsInTheDatabase() throws Throwable {
        util.put(ITEM_KEY, ITEM_ONE_KEY);
        putItem(util.get(ITEM_KEY), "payload/update_item_payload.json");
    }

    @When("^a GET HTTP request to the Item uri is done$")
    public void aGETHTTPRequestToTheItemUriIsDone() throws Throwable {
        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
        ResponseEntity<ItemData> result = client.exchange(
                (String)util.get(ENDPOINT_URI),
                HttpMethod.GET,
                new HttpEntity<>(null, headers),
                ItemData.class,
                (String)util.get(ITEM_ONE_KEY));
        util.put(RESPONSE, result);
    }

    private void putItem(String itemKey, String filename) {
        String url = String.format("%s/item/{key}", baseUrl);
        Map<String, Object> vars = new HashMap<>();
        vars.put("key", itemKey);
        ResponseEntity<Result> response = null;
        try {
            response = client.exchange(url, HttpMethod.PUT, getEntity(util.getFile(filename)), Result.class, vars);
            util.put(RESPONSE, response);
            util.remove(EXCEPTION);
        }
        catch (Exception ex) {
            util.put(EXCEPTION, ex);
        }
        if (response.getBody().isError()) {
            throw new RuntimeException(String.format("Failed to put item: %s", response.getBody().getMessage()));
        }
    }

    private Result putItemType(String itemTypeKey, String payloadFilename) {
        Result result = null;
        String url = String.format("%sitemtype/{key}", baseUrl);
        Map<String, Object> vars = new HashMap<>();
        vars.put("key", itemTypeKey);
        ResponseEntity<Result> response = null;
        try {
            response = client.exchange(url, HttpMethod.PUT, getEntity(util.getFile(payloadFilename)), Result.class, vars);
            util.put(RESPONSE, response);
            util.remove(EXCEPTION);
            result = response.getBody();
        }
        catch (Exception ex) {
            util.put(EXCEPTION, ex);
            result = new Result();
            result.setError(true);
            result.setMessage(ex.getMessage());
        }
        return result;
    }

    private Result putItemTypeAttr(String itemTypeKey, String typeAttrKey, String payloadFilename) {
        Result result = null;
        String url = String.format("%sitemtype/{item_type_key}/attribute/{attr_key}", baseUrl);
        Map<String, Object> vars = new HashMap<>();
        vars.put("item_type_key", itemTypeKey);
        vars.put("attr_key", typeAttrKey);
        ResponseEntity<Result> response = null;
        try {
            response = client.exchange(url, HttpMethod.PUT, getEntity(util.getFile(payloadFilename)), Result.class, vars);
            util.put(RESPONSE, response);
            util.remove(EXCEPTION);
            result = response.getBody();
        }
        catch (Exception ex) {
            util.put(EXCEPTION, ex);
            result = new Result();
            result.setError(true);
            result.setMessage(ex.getMessage());
        }
        return result;
    }

    private Result putLinkTypeAttr(String linkTypeKey, String typeAttrKey, String payloadFilename) {
        Result result = null;
        String url = String.format("%slinktype/{link_type_key}/attribute/{attr_key}", baseUrl);
        Map<String, Object> vars = new HashMap<>();
        vars.put("link_type_key", linkTypeKey);
        vars.put("attr_key", typeAttrKey);
        ResponseEntity<Result> response = null;
        try {
            response = client.exchange(url, HttpMethod.PUT, getEntity(util.getFile(payloadFilename)), Result.class, vars);
            util.put(RESPONSE, response);
            util.remove(EXCEPTION);
            result = response.getBody();
        }
        catch (Exception ex) {
            util.put(EXCEPTION, ex);
            result = new Result();
            result.setError(true);
            result.setMessage(ex.getMessage());
        }
        return result;
    }

    @Given("^more than one item exist in the database$")
    public void moreThanOneItemExistInTheDatabase() throws Throwable {
        theClearCMDBURLOfTheServiceIsKnown();
        aClearCMDBRequestToTheServiceIsDone();
        putModel("meta_model_1", "payload/create_model_1_payload.json");
        putItemType("item_type_1", "payload/create_item_type_1_payload.json");
        putItemTypeAttr("item_type_1", "item_type_1_COMPANY", "payload/create_item_type_attr_1_payload.json");
        putItemTypeAttr("item_type_1", "item_type_1_WBS", "payload/create_item_type_attr_2_payload.json");
        putItem("item_one", "payload/update_item_payload.json");
        putItem("item_two", "payload/update_item_payload.json");
        putItem("item_three", "payload/update_item_payload.json");
    }

    @When("^a GET HTTP request to the Item uri is done with query parameters$")
    public void aGETHTTPRequestToTheItemUriIsDoneWithQueryParameters() throws Throwable {
        DateTimeFormatter formatter = DateTimeFormatter.ofPattern("yyyyMMddHHmm");

        StringBuilder uri = new StringBuilder();
        uri.append((String)util.get(ITEM_URL));

        if (util.containsKey(CONGIG_ITEM_TYPE_ID) || util.containsKey(CONFIG_ITEM_TAG) || util.containsKey(CONFIG_ITEM_CREATED_FROM)) {
            uri.append("?");
        }

        if (util.containsKey(CONGIG_ITEM_TYPE_KEY)) {
            String typeKey = util.get(CONGIG_ITEM_TYPE_KEY);
            uri.append("type=").append(typeKey).append("&");
        }

        if (util.containsKey(CONFIG_ITEM_TAG)) {
            String tag = util.get(CONFIG_ITEM_TAG);
            uri.append("tag=").append(tag).append("&");
        }

        if (util.containsKey(CONFIG_ITEM_CREATED_FROM)) {
            ZonedDateTime from = util.get(CONFIG_ITEM_CREATED_FROM);
            uri.append("createdFrom=").append(from.format(formatter)).append("&");
        }

        if (util.containsKey(CONFIG_ITEM_CREATED_TO)) {
            ZonedDateTime to = util.get(CONFIG_ITEM_CREATED_TO);
            uri.append("createdTo=").append(to.format(formatter)).append("&");
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
        util.put(CONGIG_ITEM_TYPE_KEY, "item_type_1");
    }

    @Given("^the filtering config item tag is known$")
    public void theFilteringConfigItemTagIsKnown() throws Throwable {
        util.put(CONFIG_ITEM_TAG, "cmdb|host");
    }

    @Given("^the filtering config item date range is known$")
    public void theFilteringConfigItemDateRangeIsKnown() throws Throwable {
        util.put(CONFIG_ITEM_CREATED_FROM, ZonedDateTime.of(ZonedDateTime.now().getYear() - 100, 1, 1, 0, 0, 0, 0, ZoneId.systemDefault()));
        util.put(CONFIG_ITEM_CREATED_TO, ZonedDateTime.of(ZonedDateTime.now().getYear() + 100, 1, 1, 0, 0, 0, 0, ZoneId.systemDefault()));
    }

    @Given("^the natural key for the link is known$")
    public void theNaturalKeyForTheLinkIsKnown() throws Throwable {
        util.put(LINK_KEY, "LINK_KEY");
    }

    @When("^a PUT HTTP request with an updated JSON payload is done$")
    public void aPUTHTTPRequestWithAnUpdatedJSONPayloadIsDone() throws Throwable {
        putItem(ITEM_ONE_KEY, "payload/update_item_payload.json");
    }

    @Given("^the item type natural key is known$")
    public void theItemTypeNaturalKeyIsKnown() throws Throwable {
        util.put(CONGIG_ITEM_TYPE_KEY, "item_type_1");
    }

    private void makePutRequestWithPayload(String urlKey, String payload, String itemKey) {
        String url = util.get(urlKey);
        ResponseEntity<Result> response = null;
        try {
            response = client.exchange(url, HttpMethod.PUT, getEntity(payload), Result.class, (String) util.get(itemKey));
            util.put(RESPONSE, response);
            util.remove(EXCEPTION);
        }
        catch (Exception ex) {
            util.put(EXCEPTION, ex);
        }
    }

    private void putLink(String linkKey, String filename) {
        util.put(PAYLOAD, util.getFile(filename));
        String url = String.format("%s/link/{key}", baseUrl);
        Map<String, Object> vars = new HashMap<>();
        vars.put("key", linkKey);
        ResponseEntity<Result> response = null;
        try {
            response = client.exchange(url, HttpMethod.PUT, getEntityFromKey(PAYLOAD), Result.class, vars);
            util.put(RESPONSE, response);
            util.remove(EXCEPTION);
        }
        catch (Exception ex) {
            util.put(EXCEPTION, ex);
        }
    }

    @Given("^two items exist in the database$")
    public void twoItemsExistInTheDatabase() throws Throwable {
        putItem(ITEM_ONE_KEY, "payload/update_item_payload.json");
        putItem(ITEM_TWO_KEY, "payload/update_item_payload.json");
    }

    @Then("^the response contains (\\d+) links$")
    public void theResponseContainsLinks(int count) throws Throwable {
        ResponseEntity<LinkList> response = util.get(RESPONSE);

        LinkList links = response.getBody();
        if (links != null) {
            if (links.getValues().size() != count) {
                throw new RuntimeException(
                    String.format(
                        "Response does not contain '%s' but '%s' links.",
                        count,
                        response.getBody().getValues().size()
                    )
                );
            }
        }
        else {
            throw new RuntimeException(
                String.format(
                    "Response contains no links where '%s' were expected.",
                    count
                )
            );
        }
    }

    @Then("^the response contains more than (\\d+) items$")
    public void theResponseContainsMoreThanNumberItems(int count) {
        ResponseEntity<ItemList> response = util.get(RESPONSE);

        ItemList items = response.getBody();
        if (items != null) {
            if (items.getValues().size() <= count) {
                throw new RuntimeException(
                    String.format(
                        "Response does not contain more than '%s' items but '%s' items.",
                        count,
                        response.getBody().getValues().size()
                    )
                );
            }
        }
        else {
            throw new RuntimeException(
                String.format(
                    "Response contains no items where '%s' were expected.",
                    count
                )
            );
        }
    }

    @Then("^the reponse contains the requested item$")
    public void theReponseContainsTheRequestedItem() {
        ResponseEntity<ItemData> response = util.get(RESPONSE);
        ItemData item = response.getBody();
        if (item == null) {
            throw new RuntimeException("The response does not contain an item.");
        };
    }

    @Given("^the item type URL of the service with no query parameters exist$")
    public void theItemTypeURLOfTheServiceWithNoQueryParametersExist() {
        util.put(Key.ENDPOINT_URI, String.format("%s/itemtype", baseUrl));
    }

    @Given(("^there are item types in the database$"))
    public void thereAreItemTypesInTheDatabase(){
        for (int i = 0; i < 2; i++) {
            putItemType(
                String.format("item_type_%s", i+1) ,
                String.format("payload/create_item_type_%s_payload.json", i+1)
            );
        }
    }

    @When("^a request to GET a list of item types is done$")
    public void aRequestToGETAListOfItemTypesIsDone() {
        get(ItemTypeList.class, (String)util.get(ENDPOINT_URI));
    }

    @Then("^the response contains more than (\\d+) item types$")
    public void theResponseContainsMoreThanItemTypes(int items) {
        ResponseEntity<ItemTypeList> response = util.get(RESPONSE);
        int actual = response.getBody().getValues().size();
        if(response.getBody().getValues().size() <= items){
            throw new RuntimeException(String.format("Response contains %s item types instead of %s item types.", actual, items));
        }
    }

    @Then("^the response contains more than (\\d+) link rules$")
    public void theResponseContainsMoreThanLinkRules(int rules) {
        ResponseEntity<LinkRuleList> response = util.get(RESPONSE);
        int actual = response.getBody().getValues().size();
        if(response.getBody().getValues().size() <= rules){
            throw new RuntimeException(String.format("Response contains %s items which is less than %s items.", actual, rules));
        }
    }

    @Given("^the link between the two items exists in the database$")
    public void theLinkBetweenTheTwoItemsExistsInTheDatabase() {
        putLink(Key.LINK_ONE_KEY, "payload/create_link_payload.json");
    }

    @Given("^the item type exists in the database$")
    public void theItemTypeExistsInTheDatabase() {
        putItemType("item_type_1", "payload/create_item_type_1_payload.json");
    }

    @When("^a DELETE HTTP request with an item type key is done$")
    public void aDELETEHTTPRequestWithAnItemTypeKeyIsDone() {
        String url = (String) util.get(ITEM_TYPE_URL);
        ResponseEntity<Result> response = null;
        try {
            response = client.exchange(url, HttpMethod.DELETE, null, Result.class, (String)util.get(CONGIG_ITEM_TYPE_KEY));
            util.put(RESPONSE, response);
            util.remove(EXCEPTION);
        }
        catch (Exception ex) {
            util.put(EXCEPTION, ex);
        }
    }

    @When("^an item type DELETE HTTP request is done$")
    public void anItemTypeDELETEHTTPRequestIsDone() {
        try {
            client.delete((String) util.get(ITEM_TYPE_URL));
            util.remove(EXCEPTION);
        }
        catch (Exception ex) {
            util.put(EXCEPTION, ex);
        }
    }

    @When("^an item type PUT HTTP request with a JSON payload is done$")
    public void anItemTypePUTHTTPRequestWithAJSONPayloadIsDone() throws Throwable {
        makePutRequestWithPayload(
            ITEM_TYPE_URL,
            util.getFile("payload/create_item_type_with_meta_schema_payload.json"),
            CONGIG_ITEM_TYPE_KEY);
    }

    @Given("^the link type does not exist in the database$")
    public void theLinkTypeDoesNotExistInTheDatabase() throws Throwable {
        theClearCMDBURLOfTheServiceIsKnown();
        aClearCMDBRequestToTheServiceIsDone();
    }

    @Given("^the link type URL of the service with key is known$")
    public void theLinkTypeURLOfTheServiceWithKeyIsKnown() {
        util.put(LINK_TYPE_URL, String.format("%slinktype/{key}", baseUrl));
    }

    @Given("^the link type natural key is known$")
    public void theLinkTypeNaturalKeyIsKnown() {
        util.put(CONGIG_LINK_TYPE_KEY, "link_type_1");
    }

    @When("^a link type PUT HTTP request with a JSON payload is done$")
    public void aLinkTypePUTHTTPRequestWithAJSONPayloadIsDone() throws Throwable {
        makePutRequestWithPayload(
            LINK_TYPE_URL,
            util.getFile("payload/create_link_type_payload.json"),
            CONGIG_LINK_TYPE_KEY);
    }

    @Given("^the link type URL of the service is known$")
    public void theLinkTypeURLOfTheServiceIsKnown() {
        util.put(LINK_TYPE_URL, String.format("%slinktype", baseUrl));
    }

    @Given("^the link type exists in the database$")
    public void theLinkTypeExistsInTheDatabase() {
        theLinkTypeNaturalKeyIsKnown();
        putLinkType(util.get(CONGIG_LINK_TYPE_KEY), "payload/create_link_type_payload.json");
        putLinkTypeAttr(LINK_TYPE_ONE_KEY, "link_type_1_attr_1", "payload/create_link_type_attr_1_payload.json");
        putLinkTypeAttr(LINK_TYPE_ONE_KEY, "link_type_1_attr_2", "payload/create_link_type_attr_2_payload.json");
    }

    @When("^a DELETE HTTP request with a link type key is done$")
    public void aDELETEHTTPRequestWithALinkTypeKeyIsDone() {
        delete(util.get(LINK_TYPE_URL), util.get(CONGIG_LINK_TYPE_KEY));
    }

    @When("^a link type DELETE HTTP request is done$")
    public void aLinkTypeDELETEHTTPRequestIsDone() {
        delete(util.get(LINK_TYPE_URL), null);
    }

    @Given("^the link type URL of the service with no query parameters exist$")
    public void theLinkTypeURLOfTheServiceWithNoQueryParametersExist() {
        util.put(Key.ENDPOINT_URI, String.format("%s/linktype", baseUrl));
    }

    @When("^a request to GET a list of link types is done$")
    public void aRequestToGETAListOfLinkTypesIsDone() {
        get(LinkTypeList.class, (String)util.get(ENDPOINT_URI));
    }

    @Then("^the response contains more than (\\d+) link types$")
    public void theResponseContainsLinkTypes(int count) {
        ResponseEntity<LinkTypeList> response = util.get(RESPONSE);

        LinkTypeList links = response.getBody();
        if (links != null) {
            if (links.getValues().size() <= count) {
                throw new RuntimeException(
                    String.format(
                        "Response contains '%s' links which is less than '%s' links.",
                        response.getBody().getValues().size(),
                        count
                    )
                );
            }
        }
        else {
            throw new RuntimeException(
                    String.format(
                            "Response contains no links where '%s' were expected.",
                            count
                    )
            );
        }
    }

    @Given("^there are pre-existing Link types in the database$")
    public void thereArePreExistingLinkTypesInTheDatabase() throws Throwable {
        putData("payload/import_link_types_payload.json");
    }

    @Given("^the link rule does not exist in the database$")
    public void theLinkRuleDoesNotExistInTheDatabase() throws Throwable {
        theClearCMDBURLOfTheServiceIsKnown();
        aClearCMDBRequestToTheServiceIsDone();
    }

    @Given("^the link rule URL of the service with key is known$")
    public void theLinkRuleURLOfTheServiceWithKeyIsKnown() {
        util.put(LINK_RULE_URL, String.format("%slinkrule/{key}", baseUrl));
    }

    @Given("^the link rule natural key is known$")
    public void theLinkRuleNaturalKeyIsKnown() {
        util.put(LINK_RULE_KEY, "link_rule_1");
    }

    @When("^a link rule PUT HTTP request with a JSON payload is done$")
    public void aLinkRulePUTHTTPRequestWithAJSONPayloadIsDone() throws Throwable {
        makePutRequestWithPayload(
            LINK_RULE_URL,
            util.getFile("payload/create_link_rule_payload.json"),
            LINK_RULE_KEY);
    }

    @Given("^the link rule URL of the service is known$")
    public void theLinkRuleURLOfTheServiceIsKnown() {
        util.put(LINK_RULE_URL, String.format("%s/linkrule", baseUrl));
    }

    @When("^an link rule DELETE HTTP request is done$")
    public void anLinkRuleDELETEHTTPRequestIsDone() {
        try {
            client.delete((String) util.get(LINK_RULE_URL));
            util.remove(EXCEPTION);
        }
        catch (Exception ex) {
            util.put(EXCEPTION, ex);
        }
    }

    @Given("^the link rule exists in the database$")
    public void theLinkRuleExistsInTheDatabase() {
        putLinkRule(LINK_RULE_KEY, "payload/create_link_rule_payload.json");
    }

    @When("^a DELETE HTTP request with a link rule key is done$")
    public void aDELETEHTTPRequestWithALinkRuleKeyIsDone() {
        delete(util.get(LINK_RULE_URL), util.get(LINK_RULE_KEY));
    }

    @Given("^the link rule URL of the service with no query parameters exist$")
    public void theLinkRuleURLOfTheServiceWithNoQueryParametersExist() {
        util.put(ENDPOINT_URI, String.format("%s/linkrule", baseUrl));
    }

    @When("^a request to GET a list of link rules is done$")
    public void aRequestToGETAListOfLinkRulesIsDone() {
        get(LinkRuleList.class, (String)util.get(ENDPOINT_URI));
    }

    @Then("^the host group config items are created$")
    public void theHostGroupConfigItemsAreCreated() {
    }

    @Then("^the host config items are created$")
    public void theHostConfigItemsAreCreated() {
    }

    private void putLinkRule(String linkRuleKey, String filename) {
        util.put(PAYLOAD, util.getFile(filename));
        String url = String.format("%slinkrule/{key}", baseUrl);
        Map<String, Object> vars = new HashMap<>();
        vars.put("key", linkRuleKey);
        ResponseEntity<Result> response = null;
        try {
            response = client.exchange(url, HttpMethod.PUT, getEntityFromKey(PAYLOAD), Result.class, vars);
            util.put(RESPONSE, response);
            util.remove(EXCEPTION);
        }
        catch (Exception ex) {
            util.put(EXCEPTION, ex);
        }
    }

    /*
        a generic delete request to a specified URL in the util dictionary defined by the passed-in key
     */
    private Result delete(String url, String key) {
        Result result;
        ResponseEntity<Result> response = null;
        try {
            if (key != null) {
                response = client.exchange(url, HttpMethod.DELETE, null, Result.class, key);
            } else {
                response = client.exchange(url, HttpMethod.DELETE, null, Result.class);
            }
            util.put(RESPONSE, response);
            util.remove(EXCEPTION);
            result = response.getBody();
        }
        catch (Exception ex) {
            util.put(EXCEPTION, ex);
            result = new Result();
            result.setError(true);
            result.setMessage(ex.getMessage());
        }
        return result;
    }

    /*
        a generic get to an endpoint without parameters
     */
    private <T> void get(Class<T> cls, String url){
        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
        ResponseEntity<T> result = client.exchange(
                url,
                HttpMethod.GET,
                new HttpEntity<>(null, headers),
                cls);
        util.put(RESPONSE, result);
    }

    private void putLinkType(String linkTypeKey, String filename) {
        util.put(PAYLOAD, util.getFile(filename));
        String url = String.format("%slinktype/{key}", baseUrl);
        Map<String, Object> vars = new HashMap<>();
        vars.put("key", linkTypeKey);
        ResponseEntity<Result> response = null;
        try {
            response = client.exchange(url, HttpMethod.PUT, getEntityFromKey(PAYLOAD), Result.class, vars);
            util.put(RESPONSE, response);
            util.remove(EXCEPTION);
        }
        catch (Exception ex) {
            util.put(EXCEPTION, ex);
        }
    }

    @Given("^there are items linked to the root item in the database$")
    public void thereAreItemsLinkedToTheRootItemInTheDatabase() {
        putData("payload/create_data_payload.json");
    }

    @Given("^the URL of the tag create endpoint is known$")
    public void theURLOfTheTagCreateEndpointIsKnown() {
        util.put(TAG_CREATE_URL, String.format("%s/tag", baseUrl));
    }

    @When("^a tag creation is requested$")
    public void aTagCreationIsRequested() {
        String url = util.get(TAG_CREATE_URL);
        String payload = util.get(TAG_CREATE_PAYLOAD);
        Map<String, Object> vars = new HashMap<>();
        ResponseEntity<Result> response = null;
        try {
            response = client.exchange(url, HttpMethod.POST, getEntity(payload), Result.class, vars);
            util.put(RESPONSE, response);
            util.remove(EXCEPTION);
        }
        catch (Exception ex) {
            util.put(EXCEPTION, ex);
        }
    }

    @Given("^a payload exists with the data required to create the tag$")
    public void aPayloadExistsWithTheDataRequiredToCreateTheTag() {
        util.put(TAG_CREATE_PAYLOAD, util.getFile("payload/create_tag.json"));
    }

    @Then("^the result contains no errors$")
    public void theResultContainsNoErrors() {
        ResponseEntity<Result> result = util.get(RESPONSE);
        if(result.getBody().isError()){
            throw new RuntimeException(String.format("Result contains an error as follows: '%s'", result.getBody().getMessage()));
        };
    }

    @Given("^the URL of the tag update endpoint is known$")
    public void theURLOfTheTagUpdateEndpointIsKnown() {
        util.put(TAG_UPDATE_URL, String.format("%s/tag/{root_item_key}/{label}", baseUrl));
    }

    @Given("^the tag already exists$")
    public void theTagAlreadyExists() {
        // replays create_tag.feature
        theURLOfTheTagCreateEndpointIsKnown();
        putData("payload/create_data_payload.json");
        aPayloadExistsWithTheDataRequiredToCreateTheTag();
        aTagCreationIsRequested();
        Exception e = null;
        if (util.exists(EXCEPTION)) {
            e = util.get(EXCEPTION);
        }
        if (e != null && !e.getMessage().contains("409")) {
            throw new RuntimeException(String.format("Error creating tag"));
        };
    }

    @Given("^a payload exists with the data required to update the tag$")
    public void aPayloadExistsWithTheDataRequiredToUpdateTheTag() {
        util.put(TAG_UPDATE_PAYLOAD, util.getFile("payload/update_tag.json"));
    }

    @When("^a tag update is requested$")
    public void aTagUpdateIsRequested() {
        String url = util.get(TAG_UPDATE_URL);
        String payload = util.get(TAG_UPDATE_PAYLOAD);
        String currentLabel = util.get(TAG_LABEL);
        String itemRootKey = util.get(ROOT_ITEM_KEY);
        Map<String, Object> vars = new HashMap<>();
        vars.put("root_item_key", itemRootKey);
        vars.put("label", currentLabel);
        ResponseEntity<Result> response = null;
        try {
            response = client.exchange(url, HttpMethod.PUT, getEntity(payload), Result.class, vars);
            util.put(RESPONSE, response);
            util.remove(EXCEPTION);
        }
        catch (Exception ex) {
            util.put(EXCEPTION, ex);
        }
    }

    @Given("^the item root key of the tag is known$")
    public void theItemRootKeyOfTheTagIsKnown() {
        util.put(ROOT_ITEM_KEY, "test_inventory");
    }

    @Given("^the current label of the tag is known$")
    public void theCurrentLabelOfTheTagIsKnown() {
        util.put(TAG_LABEL, "v1");
    }

    @Given("^the URL of the tag delete endpoint is known$")
    public void theURLOfTheTagDeleteEndpointIsKnown() {
        util.put(TAG_DELETE_URL, String.format("%s/tag/{root_item_key}/{label}", baseUrl));
    }

    @When("^a tag delete is requested$")
    public void aTagDeleteIsRequested() {
        String url = util.get(TAG_DELETE_URL);
        String currentLabel = util.get(TAG_LABEL);
        String itemRootKey = util.get(ROOT_ITEM_KEY);
        Map<String, Object> vars = new HashMap<>();
        vars.put("root_item_key", itemRootKey);
        vars.put("label", currentLabel);
        ResponseEntity<Result> response = null;
        try {
            response = client.exchange(url, HttpMethod.DELETE, getEntity(), Result.class, vars);
            util.put(RESPONSE, response);
            util.remove(EXCEPTION);
        }
        catch (Exception ex) {
            util.put(EXCEPTION, ex);
        }
    }

    @Given("^the URL of the tag get endpoint is known$")
    public void theURLOfTheTagGetEndpointIsKnown() {
        util.put(TAG_LIST_URL, String.format("%s/tag/{root_item_key}", baseUrl));
    }

    @Given("^there are tags for a given item in the database$")
    public void thereAreTagsForAGivenItemInTheDatabase() {
        // for now only puts one tag in the database
        theTagAlreadyExists();
    }

    @Given("^the item root key of the tags is known$")
    public void theItemRootKeyOfTheTagsIsKnown() {
        util.put(ROOT_ITEM_KEY, "test_inventory");
    }

    @When("^a tag list for an item is requested$")
    public void aTagListForAnItemIsRequested() {
        String url = util.get(TAG_LIST_URL);
        String itemRootKey = util.get(ROOT_ITEM_KEY);
        Map<String, Object> vars = new HashMap<>();
        vars.put("root_item_key", itemRootKey);
        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
        ResponseEntity<TagList> result = client.exchange(
                url,
                HttpMethod.GET,
                new HttpEntity<>(null, headers),
                TagList.class,
                vars);
        util.put(RESPONSE, result);
    }

    @Then("^the response contains more than (\\d+) tags$")
    public void theResponseContainsMoreThanTags(int count) {
        ResponseEntity<TagList> response = util.get(RESPONSE);

        TagList items = response.getBody();
        if (items != null) {
            if (items.getValues().size() <= count) {
                throw new RuntimeException(
                    String.format(
                        "Response does not contain more than '%s' but '%s' tags.",
                        count,
                        response.getBody().getValues().size()
                    )
                );
            }
        }
        else {
            throw new RuntimeException(
                String.format(
                    "Response contains no tags where more than '%s' were expected.",
                    count
                )
            );
        }
    }

    @Given("^the tag does not already exist$")
    public void theTagDoesNotAlreadyExist() {
        // remove the tag if exists
        theURLOfTheTagDeleteEndpointIsKnown();
        theItemRootKeyOfTheTagIsKnown();
        theCurrentLabelOfTheTagIsKnown();
        aTagDeleteIsRequested();
        theResponseCodeIs(200);
    }

    @Given("^the URL of the item tree get endpoint is known$")
    public void theURLOfTheItemTreeGetEndpointIsKnown() {
        util.put(TAG_TREE_URL, String.format("%s/data/{item_key}/tag/{tag}", baseUrl));
    }

    @When("^a tag tree retrieval for the tag is requested$")
    public void aTagTreeRetrievalForTheTagIsRequested() {
        String url = util.get(TAG_TREE_URL);
        String itemRootKey = util.get(ROOT_ITEM_KEY);
        String tag = util.get(TAG_LABEL);
        Map<String, Object> vars = new HashMap<>();
        vars.put("item_key", itemRootKey);
        vars.put("tag", tag);
        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
        ResponseEntity<GraphData> result = client.exchange(
                url,
                HttpMethod.GET,
                new HttpEntity<>(null, headers),
                GraphData.class,
                vars);
        util.put(RESPONSE, result);
    }

    @Then("^the result contains the tree items and links$")
    public void theResultContainsTheTreeItemsAndLinks() {
        ResponseEntity<GraphData> response = util.get(RESPONSE);
        GraphData tree = response.getBody();
        if (tree != null) {
            if (tree.getItems().size() == 0 || tree.getLinks().size() == 0){
                throw new RuntimeException("Tree does not contain items or links.");
            }
        }
    }

    @Given("^the URL of the item tree PUT endpoint is known$")
    public void theURLOfTheItemTreePUTEndpointIsKnown() {
        util.put(IMPORT_DATA_URL, String.format("%s/data", baseUrl));
    }

    @Given("^the item tree does not exist in the database$")
    public void theItemTreeDoesNotExistInTheDatabase() {
        // do nothing for now
    }

    @Given("^a json payload with tree data exists$")
    public void aJsonPayloadWithTreeDataExists() {
        util.put(IMPORT_DATA_PAYLOAD, util.getFile("payload/create_data_payload.json"));
    }

    @When("^the creation of the tree is requested$")
    public void theCreationOfTheTreeIsRequested() {
        String payload = util.get(IMPORT_DATA_PAYLOAD);
        String url = util.get(IMPORT_DATA_URL);
        Map<String, Object> vars = new HashMap<>();
        vars.put("payload", payload);
        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
        ResponseEntity<ResultList> response = client.exchange(url, HttpMethod.PUT, getEntity(payload), ResultList.class, vars);
        util.put(RESPONSE, response);
    }

    @Then("^the result list contains no errors$")
    public void theResultListContainsNoErrors() {
        ResponseEntity<ResultList> results = util.get(RESPONSE);
        for (Result result : results.getBody().getValues()) {
            if (result.isError()) {
                throw new RuntimeException(String.format("Result contains an error as follows: '%s'", result.getMessage()));
            }
        }
    }

    @Given("^the item tree exists in the database$")
    public void theItemTreeExistsInTheDatabase() {
        theURLOfTheItemTreePUTEndpointIsKnown();
        aJsonPayloadWithTreeDataExists();
        theCreationOfTheTreeIsRequested();
        theResponseCodeIs(200);
        theResultListContainsNoErrors();
    }

    @Given("^a json payload with update tree data exists$")
    public void aJsonPayloadWithUpdateTreeDataExists() {
        util.put(UPDATE_DATA_PAYLOAD, util.getFile("payload/update_data_payload.json"));
    }

    @When("^the update of the tree is requested$")
    public void theUpdateOfTheTreeIsRequested() {
        String payload = util.get(UPDATE_DATA_PAYLOAD);
        String url = util.get(IMPORT_DATA_URL);
        Map<String, Object> vars = new HashMap<>();
        vars.put("payload", payload);
        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
        ResponseEntity<ResultList> response = client.exchange(url, HttpMethod.PUT, getEntity(payload), ResultList.class, vars);
        util.put(RESPONSE, response);
    }

    @Then("^the result list contained updated results$")
    public void theResultListContainedUpdatedResults() {
        ResponseEntity<ResultList> response = util.get(RESPONSE);
        ResultList results = response.getBody();
        boolean updateFound = false;
        for (Result result : results.getValues()) {
            updateFound = result.getOperation().equals("U");
            if (updateFound) break;
        }
        if (!updateFound){
            throw new RuntimeException("Results have not been updated");
        }
    }

    @Given("^the URL of the item tree DELETE endpoint is known$")
    public void theURLOfTheItemTreeDELETEEndpointIsKnown() {
        util.put(DELETE_TREE_URL, String.format("%s/data/{item_key}", baseUrl));
    }

    @Given("^the item key of the tree root item is known$")
    public void theItemKeyOfTheTreeRootItemIsKnown() {
        util.put(ROOT_ITEM_KEY, "test_inventory");
    }

    @When("^the deletion of the tree is requested$")
    public void theDeletionOfTheTreeIsRequested() {
        String url = util.get(DELETE_TREE_URL);
        ResponseEntity<Result> response = null;
        try {
            response = client.exchange(url, HttpMethod.DELETE, null, Result.class, (String)util.get(ROOT_ITEM_KEY));
            util.put(RESPONSE, response);
            util.remove(EXCEPTION);
        }
        catch (Exception ex) {
            util.put(EXCEPTION, ex);
        }
    }

    @Given("^there are not any tag for the root item$")
    public void thereAreNotAnyTagForTheRootItem() {
        theURLOfTheTagDeleteAllEndpointIsKnown();
        aTagDeleteAllIsRequested();
    }

    @Given("^the URL of the tag delete all endpoint is known$")
    public void theURLOfTheTagDeleteAllEndpointIsKnown() {
        util.put(TAG_DELETE_URL, String.format("%s/tag/{root_item_key}", baseUrl));
    }

    @Given("^there are more than one tags in the database$")
    public void thereAreMoreThanOneTagsInTheDatabase() {
        // do nothing for now
    }

    @When("^a tag delete all is requested$")
    public void aTagDeleteAllIsRequested() {
        deleteAllTags(util.get(TAG_DELETE_URL), (String)util.get(ROOT_ITEM_KEY));
    }

    private void deleteAllTags(String url, String key) {
        ResponseEntity<Result> response = null;
        try {
            response = client.exchange(url, HttpMethod.DELETE, null, Result.class, key);
            util.put(RESPONSE, response);
            util.remove(EXCEPTION);
        }
        catch (Exception ex) {
            util.put(EXCEPTION, ex);
        }
    }

    @Given("^the item metadata URL get by key is known$")
    public void theItemMetadataURLGetByKeyIsKnown() {
        util.put(ITEM_META_URL, String.format("%s/item/{key}/meta", baseUrl));
    }

    @When("^a GET HTTP request to the Item Metadata endpoint is done$")
    public void aGETHTTPRequestToTheItemMetadataEndpointIsDone() {
        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
        ResponseEntity<JSONObject> result = client.exchange(
                (String)util.get(ITEM_META_URL),
                HttpMethod.GET,
                new HttpEntity<>(null, headers),
                JSONObject.class,
                (String)util.get(ITEM_ONE_KEY));
        util.put(RESPONSE, result);
    }

    @Then("^the response contains the requested metadata$")
    public void theResponseContainsTheRequestedMetadata() {
        ResponseEntity<JSONObject> response = util.get(RESPONSE);
        JSONObject item = response.getBody();
        if (item == null) {
            throw new RuntimeException("Metadata not found in the response.");
        }
    }

    @Given("^the item metadata URL GET with filter is known$")
    public void theItemMetadataURLGETWithFilterIsKnown() {
        util.put(ITEM_META_URL, String.format("%s/item/{key}/meta/{filter}", baseUrl));
    }

    @Given("^an item type with filter data exists in the database$")
    public void anItemTypeWithFilterDataExistsInTheDatabase() throws Throwable {
        Result result = putItemType("item_type_with_filter","payload/create_item_type_with_filter_payload.json");
        if (result.isError()){
            throw new RuntimeException(result.getMessage());
        }
        putItemTypeAttr("item_type_with_filter", "item_type_with_filter_COMPANY", "payload/create_link_type_attr_1_payload.json");
        putItemTypeAttr("item_type_with_filter", "item_type_with_filter_WBS", "payload/create_link_type_attr_2_payload.json");
    }

    @Given("^the item with metadata exists in the database$")
    public void theItemWithMetadataExistsInTheDatabase() {
        putItem(util.get(ITEM_ONE_KEY), "payload/create_meta_test_item_payload.json");
    }

    @When("^a GET HTTP request to the Item Metadata endpoint with filter is done$")
    public void aGETHTTPRequestToTheItemMetadataEndpointWithFilterIsDone() {
        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
        ResponseEntity<JSONObject> result = client.exchange(
                (String)util.get(ITEM_META_URL),
                HttpMethod.GET,
                new HttpEntity<>(null, headers),
                JSONObject.class,
                (String)util.get(ITEM_ONE_KEY), (String) util.get(ITEM_META_FILTER));
        util.put(RESPONSE, result);
    }

    @Given("^a metadata filter key is known$")
    public void aMetadataFilterKeyIsKnown() {
        util.put(ITEM_META_FILTER, "books");
    }

    @When("^a link type GET HTTP request with the key is done$")
    public void aLinkTypeGETHTTPRequestWithTheKeyIsDone() {
        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
        ResponseEntity<JSONObject> result = client.exchange(
                (String)util.get(LINK_TYPE_URL),
                HttpMethod.GET,
                new HttpEntity<>(null, headers),
                JSONObject.class,
                (String)util.get(CONGIG_LINK_TYPE_KEY));
        util.put(RESPONSE, result);
    }

    @Then("^the response contains the link type$")
    public void theResponseContainsTheLinkType() {
        ResponseEntity<LinkTypeData> result = util.get(RESPONSE);
        if (result.getBody() == null) {
            throw new RuntimeException("The response does not contain the required link type.");
        }
    }

    @Given("^link rules exist in the database$")
    public void linkRulesExistInTheDatabase() {
        ResultList results = putData("payload/import_link_rules_payload.json");
        if (results.isError()) {
            throw new RuntimeException(results.getMessage());
        }
    }

    @Given("^the meta model does not exist in the database$")
    public void theMetaModelDoesNotExistInTheDatabase() {
        theMetaModelURLOfTheServiceWithKeyIsKnown();
        delete(util.get(MODEL_URL_WITH_KEY), util.get(META_MODEL_KEY));
    }

    @Given("^the meta model natural key is known$")
    public void theMetaModelNaturalKeyIsKnown() {
        util.put(META_MODEL_KEY, "meta_model_1");
    }

    @Given("^the meta model URL of the service with key is known$")
    public void theMetaModelURLOfTheServiceWithKeyIsKnown() {
        util.put(MODEL_URL_WITH_KEY, String.format("%s/model/{key}", baseUrl));
    }

    @When("^a meta model PUT HTTP request with a JSON payload is done$")
    public void aMetaModelPUTHTTPRequestWithAJSONPayloadIsDone() {
        util.put(RESPONSE, putModel(util.get(META_MODEL_KEY),"payload/create_model_1_payload.json"));
    }

    private ResponseEntity<Result> putModel(String key, String payload) {
        String url = String.format("%s/model/{key}", baseUrl);
        Map<String, Object> vars = new HashMap<>();
        vars.put("key", key);
        vars.put("payload", payload);
        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
        return client.exchange(url, HttpMethod.PUT, getEntity(util.getFile(payload)), Result.class, vars);
    }

    private ResponseEntity<Result> putPartition(String key, String payload) {
        String url = String.format("%s/partition/{key}", baseUrl);
        Map<String, Object> vars = new HashMap<>();
        vars.put("key", key);
        vars.put("payload", payload);
        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
        return client.exchange(url, HttpMethod.PUT, getEntity(util.getFile(payload)), Result.class, vars);
    }

    private ResponseEntity<Result> putResource(String resource, String key, String payload) {
        String url = String.format("%s/%s/{key}", baseUrl, resource);
        Map<String, Object> vars = new HashMap<>();
        vars.put("key", key);
        vars.put("payload", payload);
        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
        ResponseEntity<Result> response = client.exchange(url, HttpMethod.PUT, getEntity(util.getFile(payload)), Result.class, vars);
        util.put(RESPONSE, response);
        return response;
    }

    private ResponseEntity<Result> putPrivilege(String roleKey, String partitionKey, String payload) {
        String url = String.format("%s/privilege/{key}", baseUrl);
        Map<String, Object> vars = new HashMap<>();
        vars.put("key", String.format("%s-%s", roleKey, partitionKey));
        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
        ResponseEntity<Result> response = client.exchange(url, HttpMethod.PUT, getEntity(util.getFile(payload)), Result.class, vars);
        util.put(RESPONSE, response);
        return response;
    }

    private ResultList putData(String payloadFilePath){
        String payload = util.getFile(payloadFilePath);
        String url = String.format("%s/data", baseUrl);
        Map<String, Object> vars = new HashMap<>();
        vars.put("payload", util.getFile(payloadFilePath));
        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
        ResponseEntity<ResultList> response = client.exchange(url, HttpMethod.PUT, getEntity(payload), ResultList.class, vars);
        return response.getBody();
    }

    @Given("^the meta model exists in the database$")
    public void theMetaModelExistsInTheDatabase() {
        theMetaModelNaturalKeyIsKnown();
        putModel(util.get(META_MODEL_KEY),"payload/create_model_1_payload.json");
    }

    @When("^a meta model DELETE HTTP request with key is done$")
    public void aMetaModelDELETEHTTPRequestWithKeyIsDone() {
        delete(util.get(MODEL_URL_WITH_KEY), util.get(META_MODEL_KEY));
    }

    @Given("^the meta model URL of the service without key is known$")
    public void theMetaModelURLOfTheServiceWithoutKeyIsKnown() {
        util.put(MODEL_URL_NO_KEY, String.format("%s/model", baseUrl));
    }

    @Given("^there are a few meta models in the system$")
    public void thereAreAFewMetaModelsInTheSystem() {
        for (int i = 0; i < 2; i++) {
            putModel(
                String.format("meta_model_%s", i+1) ,
                String.format("payload/create_model_%s_payload.json", i+1)
            );
        }
    }

    @When("^a meta model GET HTTP request is done$")
    public void aMetaModelGETHTTPRequestIsDone() {
        get(ModelDataList.class, util.get(MODEL_URL_NO_KEY));
    }

    @Then("^the response contains more than (\\d+) meta models$")
    public void theResponseContainsMoreThanMetaModels(int count) {
        ResponseEntity<ModelDataList> response = util.get(RESPONSE);

        ModelDataList models = response.getBody();
        if (models != null) {
            if (models.getValues().size() <= count) {
                throw new RuntimeException(
                    String.format(
                        "Response does not contain '%s' but '%s' models.",
                        count,
                        response.getBody().getValues().size())
                );
            }
        }
        else {
            throw new RuntimeException(
                String.format(
                    "Response contains no models where '%s' were expected.",
                    count
                )
            );
        }
    }

    @When("^a meta model GET HTTP request with key is done$")
    public void aMetaModelGETHTTPRequestWithKeyIsDone() {
        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
        ResponseEntity<ModelData> result = client.exchange(
                String.format((String)util.get(MODEL_URL_WITH_KEY), baseUrl),
                HttpMethod.GET,
                new HttpEntity<>(null, headers),
                ModelData.class,
                (String)util.get(META_MODEL_KEY));
        util.put(RESPONSE, result);
    }

    @Given("^a model exists in the database$")
    public void aModelExistsInTheDatabase() {
        putModel("meta_model_1", "payload/create_model_1_payload.json");
    }

    @Given("^an item type exists in the database$")
    public void anItemTypeExistsInTheDatabase() {
        putItemType("item_type_1","payload/create_item_type_1_payload.json");
        putItemTypeAttr("item_type_1", "item_type_1_COMPANY", "payload/create_link_type_attr_1_payload.json");
        putItemTypeAttr("item_type_1", "item_type_1_WBS", "payload/create_link_type_attr_2_payload.json");
    }

    @Given("^there are not any item types associated with the model$")
    public void thereAreNotAnyItemTypesAssociatedWithTheModel() {
        Result result = delete(String.format("%s/itemtype/{key}?force=true", baseUrl), "item_type_1");
        if (result.isError()){
            throw new RuntimeException(result.getMessage());
        }
    }

    @Given("^there are not any link types associated with the model$")
    public void thereAreNotAnyLinkTypesAssociatedWithTheModel() {
        Result result = delete(String.format("%s/linktype/{key}", baseUrl), "link_type_1");
        if (result.isError()){
            throw new RuntimeException(result.getMessage());
        }
    }

    @Given("^there are not any items associated with the model$")
    public void thereAreNotAnyItemsAssociatedWithTheModel() {
        Result result = delete(String.format("%s/item", baseUrl), null);
        if (result.isError()){
            throw new RuntimeException(result.getMessage());
        }
    }

    @Given("^the link URL search with query parameters is known$")
    public void theLinkURLSearchWithQueryParametersIsKnown() {
        util.put(LINK_URL, String.format("%s/link", baseUrl));
    }

    @Given("^more than one link exist in the database$")
    public void moreThanOneLinkExistInTheDatabase() {
        putData("payload/create_data_payload.json");
    }

    @Given("^the filtering link type is known$")
    public void theFilteringLinkTypeIsKnown() {
        util.put(LINK_TYPE_FILTER, "ANSIBLE_INVENTORY");
    }

    @Given("^the filtering link tag is known$")
    public void theFilteringLinkTagIsKnown() {
        util.put(LINK_TAG_FILTER, "link|awesome");
    }

    @Given("^the filtering link date range is known$")
    public void theFilteringLinkDateRangeIsKnown() {
        util.put(LINK_CREATED_FROM_FILTER, ZonedDateTime.of(ZonedDateTime.now().getYear() - 100, 1, 1, 0, 0, 0, 0, ZoneId.systemDefault()));
        util.put(LINK_CREATED_TO_FILTER, ZonedDateTime.of(ZonedDateTime.now().getYear() + 100, 1, 1, 0, 0, 0, 0, ZoneId.systemDefault()));
    }

    @Then("^the response contains more than (\\d+) links$")
    public void theResponseContainsMoreThanLinks(int count) {
        ResponseEntity<LinkList> response = util.get(RESPONSE);

        LinkList links = response.getBody();
        if (links != null) {
            if (links.getValues().size() <= count) {
                throw new RuntimeException(
                    String.format(
                        "Response does not contain '%s' but '%s' links.",
                        count,
                        response.getBody().getValues().size())
                );
            }
        }
        else {
            throw new RuntimeException(
                String.format(
                    "Response contains no links where '%s' were expected.",
                    count
                )
            );
        }
    }

    @When("^a GET HTTP request to the Link uri is done with query parameters$")
    public void aGETHTTPRequestToTheLinkUriIsDoneWithQueryParameters() {
        DateTimeFormatter formatter = DateTimeFormatter.ofPattern("yyyyMMddHHmm");

        StringBuilder uri = new StringBuilder();
        uri.append((String)util.get(LINK_URL));

        if (util.containsKey(LINK_TAG_FILTER) || util.containsKey(LINK_TYPE_FILTER) || util.containsKey(LINK_CREATED_FROM_FILTER)) {
            uri.append("?");
        }

        if (util.containsKey(LINK_TYPE_FILTER)) {
            String typeKey = util.get(LINK_TYPE_FILTER);
            uri.append("type=").append(typeKey).append("&");
        }

        if (util.containsKey(LINK_TAG_FILTER)) {
            String tag = util.get(LINK_TAG_FILTER);
            uri.append("tag=").append(tag).append("&");
        }

        if (util.containsKey(LINK_CREATED_FROM_FILTER)) {
            ZonedDateTime from = util.get(LINK_CREATED_FROM_FILTER);
            uri.append("createdFrom=").append(from.format(formatter)).append("&");
        }

        if (util.containsKey(LINK_CREATED_TO_FILTER)) {
            ZonedDateTime to = util.get(LINK_CREATED_TO_FILTER);
            uri.append("createdTo=").append(to.format(formatter)).append("&");
        }

        String uriString = uri.toString();
        if (uriString.endsWith("&")) {
            uriString = uriString.substring(0, uriString.length() - 1);
        }

        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
        ResponseEntity<LinkList> result = client.exchange(
                uriString,
                HttpMethod.GET,
                new HttpEntity<>(null, headers),
                LinkList.class);
        util.put(RESPONSE, result);
    }

    @Given("^an item link natural key is known$")
    public void anItemLinkNaturalKeyIsKnown() {
        util.put(LINK_KEY, "link_type_1");
    }

    @Given("^the link URL search by key is known$")
    public void theLinkURLSearchByKeyIsKnown() {
        util.put(LINK_WITH_KEY_URL, String.format("%s/link/{key}", baseUrl));
    }

    @Given("^the link exists in the database$")
    public void theLinkExistsInTheDatabase() {
        putLink(util.get(LINK_KEY), "payload/create_link_payload.json");
    }

    @When("^a GET HTTP request to the Link with Key URL is done$")
    public void aGETHTTPRequestToTheLinkWithKeyURLIsDone() {
        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
        ResponseEntity<LinkData> result = client.exchange(
                String.format((String)util.get(LINK_WITH_KEY_URL), baseUrl),
                HttpMethod.GET,
                new HttpEntity<>(null, headers),
                LinkData.class,
                (String)util.get(LINK_KEY));
        util.put(RESPONSE, result);
    }

    @Then("^the reponse contains the requested link$")
    public void theReponseContainsTheRequestedLink() {
        ResponseEntity<LinkData> response = util.get(RESPONSE);
        LinkData link = response.getBody();
        if (link == null){
            throw new RuntimeException("The response does not contain a link.");
        };
    }

    @Given("^the live URL of the service is known$")
    public void theLiveURLOfTheServiceIsKnown() {
        util.put(LIVE_URL, String.format("%slive", baseUrl));
    }

    @Given("^the item types to and from exists in the database$")
    public void theItemTypesToAndFromExistsInTheDatabase() {
        putItemType("item_type_1", "payload/create_item_type_1_payload.json");
        putItemType("item_type_2", "payload/create_item_type_2_payload.json");
    }

    @Given("^a partition natural key is known$")
    public void aPartitionNaturalKeyIsKnown() {
        util.put(PARTITION_ONE_KEY, PARTITION_ONE_KEY);
    }

    @Given("^the partition does not exist in the database$")
    public void thePartitionDoesNotExistInTheDatabase() {
        Result r = delete(String.format("%spartition/{partition_key}", baseUrl), util.get(PARTITION_ONE_KEY));
        if (r.isError()) {
            throw new RuntimeException(r.getMessage());
        }
    }

    @Given("^the partition PUT URL by key is known$")
    public void thePartitionPUTURLByKeyIsKnown() {
        util.put(ENDPOINT_URI, String.format("%spartition/{partition_key}", baseUrl));
    }

    @Given("^the partition exists in the database$")
    public void thePartitionExistsInTheDatabase() {
        putPartition(PARTITION_ONE_KEY, "payload/create_partition_1_payload.json");
    }

    @Given("^the partition DELETE URL by key is known$")
    public void thePartitionDELETEURLByKeyIsKnown() {
        util.put(PARTITION_DELETE_URL, String.format("%spartition/{partition_key}", baseUrl));
    }

    @When("^a PUT HTTP request to the partition endpoint with a new JSON payload is done$")
    public void aPUTHTTPRequestToThePartitionEndpointWithANewJSONPayloadIsDone() {
        ResponseEntity<Result> response = putPartition(PARTITION_ONE_KEY, "payload/create_partition_1_payload.json");
        util.put(RESPONSE, response);
    }

    @When("^a DELETE HTTP request to the partition resource with an item key is done$")
    public void aDELETEHTTPRequestToThePartitionResourceWithAnItemKeyIsDone() {
        delete(util.get(PARTITION_DELETE_URL), util.get(PARTITION_ONE_KEY));
    }

    @Given("^there are multiple partitions in the database$")
    public void thereAreMultiplePartitionsInTheDatabase() {
        delete( String.format("%spartition/{partition_key}", baseUrl), "PART_01");
        delete( String.format("%spartition/{partition_key}", baseUrl), "PART_02");
        delete( String.format("%spartition/{partition_key}", baseUrl), "PART_03");
        putPartition("PART_01", "payload/create_partition_1_payload.json");
        putPartition("PART_02", "payload/create_partition_2_payload.json");
        putPartition("PART_03", "payload/create_partition_3_payload.json");
    }

    @Given("^the partition URL of the service with no query parameters exist$")
    public void thePartitionURLOfTheServiceWithNoQueryParametersExist() {
        util.put(PARTITION_GET_URL, String.format("%spartition", baseUrl));
    }

    @When("^a request to GET a list of partitions is done$")
    public void aRequestToGETAListOfPartitionsIsDone() {
        get(PartitionDataList.class, util.get(PARTITION_GET_URL));
    }

    @Then("^the response contains more than (\\d+) partitions$")
    public void theResponseContainsMoreThanPartitions(int count) {
        ResponseEntity<PartitionDataList> response = util.get(RESPONSE);

        PartitionDataList items = response.getBody();
        if (items != null) {
            if (items.getValues().size() <= count) {
                throw new RuntimeException(
                    String.format(
                        "Response does not contain more than '%s' items but '%s' items.",
                        count,
                        response.getBody().getValues().size()
                    )
                );
            }
        }
        else {
            throw new RuntimeException(
                    String.format(
                            "Response contains no items where '%s' were expected.",
                            count
                    )
            );
        }
    }

    @Given("^the partition natural key is known$")
    public void thePartitionNaturalKeyIsKnown() {
        util.put(PARTITION_ONE_KEY, PARTITION_ONE_KEY);
    }

    @Given("^the partition is in the database$")
    public void thePartitionIsInTheDatabase() {
        putPartition(util.get(PARTITION_ONE_KEY), "payload/create_partition_1_payload.json");
    }

    @When("^a request to GET the partition is made$")
    public void aRequestToGETThePartitionIsMade() {
        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
        ResponseEntity<PartitionData> result = client.exchange(
                String.format("%spartition/{key}", baseUrl),
                HttpMethod.GET,
                new HttpEntity<>(null, headers),
                PartitionData.class,
                (String)util.get(PARTITION_ONE_KEY));
        util.put(RESPONSE, result);
    }

    @Then("^the response contains the requested partition$")
    public void theResponseContainsTheRequestedPartition() {
        ResponseEntity<PartitionData> partData = util.get(RESPONSE);
        PartitionData part = partData.getBody();
        if (part == null) {
            throw new RuntimeException("Partition not found");
        }
    }

    @Given("^a role natural key is known$")
    public void aRoleNaturalKeyIsKnown() {
        util.put(ROLE_ONE_KEY, ROLE_ONE_KEY);
    }

    @Given("^the role exists in the database$")
    public void theRoleExistsInTheDatabase() {
        aRoleNaturalKeyIsKnown();
        aDELETEHTTPRequestToTheRoleResourceWithAnItemKeyIsDone();
        ResponseEntity<Result> response = putResource("role", util.get(ROLE_ONE_KEY), "payload/create_role_1_payload.json");
        Result result = response.getBody();
    }

    @Given("^the role DELETE URL by key is known$")
    public void theRoleDELETEURLByKeyIsKnown() {
        util.put(ROLE_DELETE_URL, String.format("%srole/{key}", baseUrl));
    }

    @When("^a DELETE HTTP request to the role resource with an item key is done$")
    public void aDELETEHTTPRequestToTheRoleResourceWithAnItemKeyIsDone() {
        theRoleDELETEURLByKeyIsKnown();
        delete(util.get(ROLE_DELETE_URL), util.get(ROLE_ONE_KEY));
    }

    @Given("^the role does not exist in the database$")
    public void theRoleDoesNotExistInTheDatabase() {
        aDELETEHTTPRequestToTheRoleResourceWithAnItemKeyIsDone();
    }

    @Given("^the role PUT URL by key is known$")
    public void theRolePUTURLByKeyIsKnown() {
        util.put(ROLE_PUT_URL, String.format("%srole/{key}", baseUrl));
    }

    @When("^a PUT HTTP request to the role endpoint with a new JSON payload is done$")
    public void aPUTHTTPRequestToTheRoleEndpointWithANewJSONPayloadIsDone() {
        putResource("role", util.get(ROLE_ONE_KEY), "payload/create_role_1_payload.json");
    }

    @When("^a request to GET the role is made$")
    public void aRequestToGETTheRoleIsMade() {
        getByKey("role", util.get(ROLE_ONE_KEY));
    }

    private void getByKey(String resource, String key) {
        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
        ResponseEntity<PartitionData> result = client.exchange(
                String.format("%s%s/{key}", baseUrl, resource),
                HttpMethod.GET,
                new HttpEntity<>(null, headers),
                PartitionData.class,
                key);
        util.put(RESPONSE, result);
    }

    @Given("^there are multiple roles in the database$")
    public void thereAreMultipleRolesInTheDatabase() {
        delete(String.format("%srole/{key}", baseUrl), "ROLE_01");
        delete(String.format("%srole/{key}", baseUrl), "ROLE_02");
        delete(String.format("%srole/{key}", baseUrl), "ROLE_03");
        putPartition("ROLE_01", "payload/create_role_1_payload.json");
        putPartition("ROLE_02", "payload/create_role_2_payload.json");
        putPartition("ROLE_03", "payload/create_role_3_payload.json");
    }

    @Given("^the role URL of the service with no query parameters exist$")
    public void theRoleURLOfTheServiceWithNoQueryParametersExist() {
        util.put(ROLE_GET_URL, String.format("%srole", baseUrl));
    }

    @When("^a request to GET a list of roles is done$")
    public void aRequestToGETAListOfRolesIsDone() {
        get(RoleDataList.class, util.get(ROLE_GET_URL));
    }

    @Then("^the response contains more than (\\d+) roles$")
    public void theResponseContainsMoreThanRoles(int count) {
        ResponseEntity<RoleDataList> response = util.get(RESPONSE);

        RoleDataList items = response.getBody();
        if (items != null) {
            if (items.getValues().size() <= count) {
                throw new RuntimeException(
                    String.format(
                        "Response does not contain more than '%s' roles but '%s' roles.",
                        count,
                        response.getBody().getValues().size()
                    )
                );
            }
        }
        else {
            throw new RuntimeException(
                String.format(
                    "Response contains no roles where '%s' were expected.",
                    count
                )
            );
        }
    }

    @When("^a PUT HTTP request to the privilege endpoint with a new JSON payload is done$")
    public void aPUTHTTPRequestToThePrivilegeEndpointWithANewJSONPayloadIsDone() {
        putPrivilege(ROLE_ONE_KEY, PARTITION_ONE_KEY,"payload/create_privilege_1_payload.json");
    }

    @Given("^the privilege does not exist in the database$")
    public void thePrivilegeDoesNotExistInTheDatabase() {
        aDELETEHTTPRequestToThePrivilegeEndpointIsDone();
    }

    @Given("^the privilege exists in the database$")
    public void thePrivilegeExistsInTheDatabase() {
        putPrivilege(ROLE_ONE_KEY, PARTITION_ONE_KEY,"payload/create_privilege_1_payload.json");
    }

    @When("^a DELETE HTTP request to the privilege endpoint is done$")
    public void aDELETEHTTPRequestToThePrivilegeEndpointIsDone() {
        deletePrivilege(ROLE_ONE_KEY, PARTITION_ONE_KEY);
    }

    @Given("^there are multiple privileges for a role in the database$")
    public void thereAreMultiplePrivilegesForARoleInTheDatabase() {
        delete(String.format("%s/partition/{key}", baseUrl), "PART_01");
        delete(String.format("%s/role/{key}", baseUrl), "ROLE_01");
        deletePrivilege("ROLE_01", "PART_01");
        putResource("partition", "PART_01", "payload/create_partition_1_payload.json");
        putResource("partition", "PART_02", "payload/create_partition_2_payload.json");
        putResource("role", "ROLE_01", "payload/create_role_1_payload.json");
        putPrivilege("ROLE_01", "PART_01","payload/create_privilege_1_payload.json");
        putPrivilege("ROLE_01", "PART_02","payload/create_privilege_2_payload.json");
    }

    @When("^a request to GET a list of privileges by role is done$")
    public void aRequestToGETAListOfPrivilegesByRoleIsDone() {
        get(PrivilegeDataList.class, String.format("%s/role/ROLE_01/privilege", baseUrl));
    }

    @Then("^the response contains more than (\\d+) privileges$")
    public void theResponseContainsMoreThanPrivileges(int arg0) {
    }

    private void deletePrivilege(String roleKey, String partitionKey) {
        Result result;
        ResponseEntity<Result> response = null;
        response = client.exchange(
                String.format("%s/privilege/{key}", baseUrl),
                HttpMethod.DELETE,
                null,
                Result.class,
                String.format("%s-%s", ROLE_ONE_KEY, PARTITION_ONE_KEY));
        util.put(RESPONSE, response);
        util.remove(EXCEPTION);
        result = response.getBody();
    }

    @When("^the readiness probe is checked$")
    public void theReadinessProbeIsChecked() {
        String url = baseUrl + "/ready";
        ResponseEntity<Result> response = null;
        try {
            response = client.exchange(url, HttpMethod.GET, null, Result.class);
            util.put(RESPONSE, response);
            util.remove(EXCEPTION);
        }
        catch (Exception ex) {
            util.put(EXCEPTION, ex);
        }
    }

    @Given("^the database does not exist$")
    public void theDatabaseDoesNotExist() {
        // delete the database here
    }

    @Then("^the database is deployed$")
    public void theDatabaseIsDeployed() {
        if (util.containsKey(EXCEPTION)){
            throw new RuntimeException((Exception)util.get(EXCEPTION));
        }
    }

    @Given("^the item type attribute natural key is known$")
    public void theItemTypeAttributeNaturalKeyIsKnown() {
        util.put(ITEM_TYPE_ATTR_ONE_KEY, ITEM_TYPE_ATTR_ONE_KEY);
    }

    @When("^a PUT HTTP request with a JSON payload is done for an attribute of an item type$")
    public void aPUTHTTPRequestWithAJSONPayloadIsDoneForAnAttributeOfAnItemType() {
        putItemTypeAttr("item_type_1", util.get(ITEM_TYPE_ATTR_ONE_KEY), "payload/create_item_type_attr_1_payload.json");
    }

    @Given("^the item type attribute does not exist in the database$")
    public void theItemTypeAttributeDoesNotExistInTheDatabase() {
        Result r = delete(String.format("%sitemtype/%s/attribute/%s", baseUrl, util.get(ITEM_TYPE_ONE_KEY), util.get(ITEM_TYPE_ATTR_ONE_KEY))+"?force=true", "item_type_1");
        if (r.isError()) {
            throw new RuntimeException(r.getMessage());
        }
    }

    @Given("^the key of the item type is known$")
    public void theKeyOfTheItemTypeIsKnown() {
        util.put(ITEM_TYPE_ONE_KEY, ITEM_TYPE_ONE_KEY);
    }

    @Given("^the key of the type attribute for the item type is known$")
    public void theKeyOfTheTypeAttributeForTheItemTypeIsKnown() {
        util.put(ITEM_TYPE_ATTR_ONE_KEY, ITEM_TYPE_ATTR_ONE_KEY);
    }

    @Given("^there are item type attributes for the item types in the database$")
    public void thereAreItemTypeAttributesForTheItemTypesInTheDatabase() {
        putItemTypeAttr("item_type_1", "item_type_1_attr_1", "payload/create_item_type_attr_1_payload.json");
        putItemTypeAttr("item_type_1", "item_type_1_attr_2", "payload/create_item_type_attr_1_payload.json");
        putItemTypeAttr("item_type_1", "item_type_1_attr_3", "payload/create_item_type_attr_1_payload.json");
    }

    @Given("^the item type attribute URL exist$")
    public void theItemTypeAttributeURLExist() {
        util.put(ENDPOINT_URI, String.format("%sitemtype/{item_type_key}/attribute", baseUrl));
    }

    @When("^a request to GET a list of item type attributes is done$")
    public void aRequestToGETAListOfItemTypeAttributesIsDone() {
        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
        Map<String, Object> vars = new HashMap<>();
        vars.put("item_type_key", ITEM_TYPE_ONE_KEY);
        ResponseEntity<TypeAttrList> result = client.exchange(
                util.get(ENDPOINT_URI),
                HttpMethod.GET,
                new HttpEntity<>(null, headers),
                TypeAttrList.class,
                vars);
        util.put(RESPONSE, result);
    }

    @Then("^the response contains more than (\\d+) item type attributes$")
    public void theResponseContainsMoreThanItemTypeAttributes(int items) {
        ResponseEntity<TypeAttrList> results = (ResponseEntity<TypeAttrList>) util.get(RESPONSE);
        if (results.getBody().getValues().size() <= items) {
            throw new RuntimeException("Not enough items in result");
        }
    }

    @Given("^the key of the link type is known$")
    public void theKeyOfTheLinkTypeIsKnown() {
        util.put(LINK_TYPE_ONE_KEY, LINK_TYPE_ONE_KEY);
    }

    @Given("^the key of the type attribute for the link type is known$")
    public void theKeyOfTheTypeAttributeForTheLinkTypeIsKnown() {
        util.put(LINK_TYPE_ATTR_ONE_KEY, LINK_TYPE_ATTR_ONE_KEY);
    }

    @Given("^the link type attribute does not exist in the database$")
    public void theLinkTypeAttributeDoesNotExistInTheDatabase() {
        Result r = delete(String.format("%slinktype/%s/attribute/%s", baseUrl, util.get(LINK_TYPE_ONE_KEY), util.get(LINK_TYPE_ATTR_ONE_KEY))+"?force=true", "link_type_1");
        if (r.isError()){
            throw new RuntimeException(r.getMessage());
        }
    }

    @Given("^an link type exists in the database$")
    public void anLinkTypeExistsInTheDatabase() {
        putLinkType(util.get(LINK_TYPE_ONE_KEY),"payload/create_link_type_1_payload.json");
    }

    @When("^a PUT HTTP request with a JSON payload is done for an attribute of a link type$")
    public void aPUTHTTPRequestWithAJSONPayloadIsDoneForAnAttributeOfALinkType() {
        putLinkTypeAttr(util.get(LINK_TYPE_ONE_KEY), util.get(LINK_TYPE_ATTR_ONE_KEY), "payload/create_link_type_attr_1_payload.json");
    }

    @Given("^the link type attribute natural key is known$")
    public void theLinkTypeAttributeNaturalKeyIsKnown() {
        util.put(LINK_TYPE_ATTR_ONE_KEY, LINK_TYPE_ATTR_ONE_KEY);
    }

    @Given("^there are link types in the database$")
    public void thereAreLinkTypesInTheDatabase() {
        for (int i = 0; i < 2; i++) {
            putItemType(
                String.format("link_type_%s", i+1) ,
                String.format("payload/create_link_type_%s_payload.json", i+1)
            );
        }
    }

    @Given("^there are link type attributes for the link types in the database$")
    public void thereAreLinkTypeAttributesForTheLinkTypesInTheDatabase() {
        putLinkTypeAttr(LINK_TYPE_ONE_KEY, "link_type_1_attr_1", "payload/create_link_type_attr_1_payload.json");
        putLinkTypeAttr(LINK_TYPE_ONE_KEY, "link_type_1_attr_2", "payload/create_link_type_attr_2_payload.json");
    }

    @Given("^the link type attribute URL exist$")
    public void theLinkTypeAttributeURLExist() {
        util.put(ENDPOINT_URI, String.format("%slinktype/{link_type_key}/attribute", baseUrl));
    }

    @When("^a request to GET a list of link type attributes is done$")
    public void aRequestToGETAListOfLinkTypeAttributesIsDone() {
        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
        Map<String, Object> vars = new HashMap<>();
        vars.put("link_type_key", "link_type_1");
        ResponseEntity<TypeAttrList> result = client.exchange(
                util.get(ENDPOINT_URI),
                HttpMethod.GET,
                new HttpEntity<>(null, headers),
                TypeAttrList.class,
                vars);
        util.put(RESPONSE, result);
    }
}
