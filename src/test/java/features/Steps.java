package features;

import cucumber.api.java.en.And;
import cucumber.api.java.en.Given;
import cucumber.api.java.en.Then;
import cucumber.api.java.en.When;
import org.gatblau.onix.Info;
import org.gatblau.onix.data.*;
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
        util.put(ITEM_URL, String.format("%sitem/{key}", baseUrl));
    }

    @Given("^the item URL search with query parameters is known$")
    public void theItemURLSearchWithQueryParametersIsKnown() throws Throwable {
        util.put(ITEM_URL, String.format("%s/item", baseUrl));
    }

    @And("^the response code is (\\d+)$")
    public void theResponseCodeIs(int responseCode)  {
        if (util.containsKey(EXCEPTION)) {
            RuntimeException ex = util.get(EXCEPTION);
            throw ex;
        }
        ResponseEntity<Result> response = util.get(RESPONSE);
        if (response.getStatusCodeValue() != responseCode) {
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
        putItem(ITEM_ONE_KEY, "payload/create_item_payload.json");
    }

    @And("^a configuration item natural key is known$")
    public void aConfigurationItemNaturalKeyIsKnown() throws Throwable {
        util.put(ITEM_ONE_KEY, ITEM_ONE_KEY);
    }

    @And("^the service responds with action \"([^\"]*)\"$")
    public void theServiceRespondsWithAction(String action) throws Throwable {
        ResponseEntity<Result> response = util.get(RESPONSE);
        Result result = response.getBody();
        assert (result.getOperation().equals(action));
    }

    @And("^the item exist in the database$")
    public void theItemExistInTheDatabase() throws Throwable {
        theItemDoesNotExistInTheDatabase();
        theItemURLSearchByKeyIsKnown();
        aJsonPayloadWithNewItemInformationExists();
        aPUTTHTTPRequestWithANewJSONPayloadIsDone();
    }

    @Given("^the item type does not exist in the database$")
    public void theItemTypeDoesNotExistInTheDatabase() throws Throwable {
        theItemTypeURLOfTheServiceIsKnown();
        thereIsNotAnyErrorInTheResponse();
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
        delete(ITEM_URL, ITEM_KEY);
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
        util.put(PAYLOAD, util.getFile(filename));
        String url = String.format("%s/item/{key}", baseUrl);
        Map<String, Object> vars = new HashMap<>();
        vars.put("key", itemKey);
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

    private void putItemType(String itemTypeKey, String filename) {
        util.put(PAYLOAD, util.getFile(filename));
        String url = String.format("%sitemtype/{key}", baseUrl);
        Map<String, Object> vars = new HashMap<>();
        vars.put("key", itemTypeKey);
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

    @Given("^more than one item exist in the database$")
    public void moreThanOneItemExistInTheDatabase() throws Throwable {
        theClearCMDBURLOfTheServiceIsKnown();
        aClearCMDBRequestToTheServiceIsDone();
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

        if (util.containsKey(CONGIG_ITEM_TYPE_KEY)) {
            String typeKey = util.get(CONGIG_ITEM_TYPE_KEY);
            uri.append("type=").append(typeKey).append("&");
        }

        if (util.containsKey(CONFIG_ITEM_TAG)) {
            String tag = util.get(CONFIG_ITEM_TAG);
            uri.append("tag=").append(tag).append("&");
        }

        if (util.containsKey(CONFIG_ITEM_UPDATED_FROM)) {
            ZonedDateTime from = util.get(CONFIG_ITEM_UPDATED_FROM);
            uri.append("createdFrom=").append(from.format(formatter)).append("&");
        }

        if (util.containsKey(CONFIG_ITEM_UPDATED_TO)) {
            ZonedDateTime to = util.get(CONFIG_ITEM_UPDATED_TO);
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
        util.put(CONGIG_ITEM_TYPE_KEY, "HOST");
    }

    @Given("^the filtering config item tag is known$")
    public void theFilteringConfigItemTagIsKnown() throws Throwable {
        util.put(CONFIG_ITEM_TAG, "cmdb|host");
    }

    @Given("^the filtering config item date range is known$")
    public void theFilteringConfigItemDateRangeIsKnown() throws Throwable {
        util.put(CONFIG_ITEM_UPDATED_FROM, ZonedDateTime.of(ZonedDateTime.now().getYear() - 100, 1, 1, 0, 0, 0, 0, ZoneId.systemDefault()));
        util.put(CONFIG_ITEM_UPDATED_TO, ZonedDateTime.of(ZonedDateTime.now().getYear() + 100, 1, 1, 0, 0, 0, 0, ZoneId.systemDefault()));
    }

    @Given("^the natural key for the link is known$")
    public void theNaturalKeyForTheLinkIsKnown() throws Throwable {
        util.put(LINK_KEY, "LINK_KEY");
    }

    @When("^a PUT HTTP request with an updated JSON payload is done$")
    public void aPUTHTTPRequestWithAnUpdatedJSONPayloadIsDone() throws Throwable {
        putItem(ITEM_ONE_KEY, "payload/update_item_payload.json");
    }

    @Given("^a json payload with new item type information exists$")
    public void aJsonPayloadWithNewItemTypeInformationExists() throws Throwable {
        util.put(PAYLOAD, util.getFile("payload/create_item_type_payload.json"));
    }

    @Given("^the item type natural key is known$")
    public void theItemTypeNaturalKeyIsKnown() throws Throwable {
        util.put(CONGIG_ITEM_TYPE_KEY, "item_type_1");
    }

//    @When("^a PUT HTTP request with a JSON payload is done$")
    private void makePutRequestWithPayload(String urlKey, String payloadKey, String itemKey) throws Throwable {
        String url = util.get(urlKey);
        ResponseEntity<Result> response = null;
        try {
            response = client.exchange(url, HttpMethod.PUT, getEntityFromKey(payloadKey), Result.class, (String) util.get(itemKey));
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

//    @Given("^two links between the two configuration items exist in the database$")
//    public void twoLinksBetweenTheTwoConfigurationItemsExistInTheDatabase() throws Throwable {
//        putLink(Key.LINK_ONE_KEY, "payload/create_link_payload.json");
//        putLink(Key.LINK_TWO_KEY, "payload/create_link_payload.json");
//    }

    @When("^a GET HTTP request to the Link by Item resource is done$")
    public void aGETHTTPRequestToTheLinkByItemResourceIsDone() throws Throwable {
        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
        ResponseEntity<LinkList> result = client.exchange(
            String.format((String)util.get(ENDPOINT_URI), baseUrl),
            HttpMethod.GET,
            new HttpEntity<>(null, headers),
            LinkList.class,
            (String)util.get(ITEM_ONE_KEY));
        util.put(RESPONSE, result);
    }

    @Then("^the response contains (\\d+) links$")
    public void theResponseContainsLinks(int count) throws Throwable {
        ResponseEntity<LinkList> response = util.get(RESPONSE);

        LinkList links = response.getBody();
        if (links != null) {
            if (links.getItems().size() != count) {
                throw new RuntimeException(
                    String.format(
                        "Response does not contain '%s' but '%s' links.",
                        count,
                        response.getBody().getItems().size()
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

    @Then("^the response contains (\\d+) items$")
    public void theResponseContainsNumberItems(int count) {
        ResponseEntity<ItemList> response = util.get(RESPONSE);

        ItemList items = response.getBody();
        if (items != null) {
            if (items.getItems().size() != count) {
                throw new RuntimeException(
                    String.format(
                        "Response does not contain '%s' but '%s' items.",
                        count,
                        response.getBody().getItems().size()
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
        assert(item != null);
    }

    @Given("^the item type URL of the service with no query parameters exist$")
    public void theItemTypeURLOfTheServiceWithNoQueryParametersExist() {
        util.put(Key.ENDPOINT_URI, String.format("%s/itemtype", baseUrl));
    }

    @When("^a request to GET a list of item types is done$")
    public void aRequestToGETAListOfItemTypesIsDone() {
        get(ItemTypeList.class);
    }

    @Then("^the response contains (\\d+) item types$")
    public void theResponseContainsItemTypes(int items) {
        ResponseEntity<ItemTypeList> response = util.get(RESPONSE);
        int actual = response.getBody().getItems().size();
        if(response.getBody().getItems().size() != items){
            throw new RuntimeException(String.format("Response contains %s items instead of %s items.", actual, items));
        }
    }

    @Then("^the response contains more than (\\d+) link rules$")
    public void theResponseContainsMoreThanLinkRules(int rules) {
        ResponseEntity<LinkRuleList> response = util.get(RESPONSE);
        int actual = response.getBody().getItems().size();
        if(response.getBody().getItems().size() <= rules){
            throw new RuntimeException(String.format("Response contains %s items which is less than %s items.", rules, actual));
        }
    }

    @Given("^the link between the two items exists in the database$")
    public void theLinkBetweenTheTwoItemsExistsInTheDatabase() {
        putLink(Key.LINK_ONE_KEY, "payload/create_link_payload.json");
    }

    @Given("^the item type exists in the database$")
    public void theItemTypeExistsInTheDatabase() {
        putItemType(util.get(CONGIG_ITEM_TYPE_KEY), "payload/create_item_type_payload.json");
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
        makePutRequestWithPayload(ITEM_TYPE_URL, PAYLOAD, CONGIG_ITEM_TYPE_KEY);
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

    @Given("^a json payload with new link type information exists$")
    public void aJsonPayloadWithNewLinkTypeInformationExists() {
        String payload = util.getFile("payload/create_link_type_payload.json");
        util.put(Key.PAYLOAD, payload);
    }

    @When("^a link type PUT HTTP request with a JSON payload is done$")
    public void aLinkTypePUTHTTPRequestWithAJSONPayloadIsDone() throws Throwable {
        makePutRequestWithPayload(LINK_TYPE_URL, PAYLOAD, CONGIG_LINK_TYPE_KEY);
    }

    @Given("^the link type URL of the service is known$")
    public void theLinkTypeURLOfTheServiceIsKnown() {
        util.put(LINK_TYPE_URL, String.format("%slinktype", baseUrl));
    }

    @Given("^the link type exists in the database$")
    public void theLinkTypeExistsInTheDatabase() {
        putLinkType(util.get(CONGIG_LINK_TYPE_KEY), "payload/create_link_type_payload.json");
    }

    @When("^a DELETE HTTP request with a link type key is done$")
    public void aDELETEHTTPRequestWithALinkTypeKeyIsDone() {
        delete(LINK_TYPE_URL, CONGIG_LINK_TYPE_KEY);
    }

    @When("^a link type DELETE HTTP request is done$")
    public void aLinkTypeDELETEHTTPRequestIsDone() {
        delete(LINK_TYPE_URL, null);
    }

    @Given("^the link type URL of the service with no query parameters exist$")
    public void theLinkTypeURLOfTheServiceWithNoQueryParametersExist() {
        util.put(Key.ENDPOINT_URI, String.format("%s/linktype", baseUrl));
    }

    @When("^a request to GET a list of link types is done$")
    public void aRequestToGETAListOfLinkTypesIsDone() {
        get(LinkTypeList.class);
    }

    @Then("^the response contains more than (\\d+) link types$")
    public void theResponseContainsLinkTypes(int count) {
        ResponseEntity<LinkTypeList> response = util.get(RESPONSE);

        LinkTypeList links = response.getBody();
        if (links != null) {
            if (links.getItems().size() <= count) {
                throw new RuntimeException(
                    String.format(
                        "Response contains '%s' links which is less than '%s' links.",
                        response.getBody().getItems().size(),
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
        // there are already 3 pre-existing system item types in the database so do not do anything
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

    @Given("^a json payload with new link rule information exists$")
    public void aJsonPayloadWithNewLinkRuleInformationExists() {
        String payload = util.getFile("payload/create_link_rule_payload.json");
        util.put(Key.PAYLOAD, payload);
    }

    @When("^a link rule PUT HTTP request with a JSON payload is done$")
    public void aLinkRulePUTHTTPRequestWithAJSONPayloadIsDone() throws Throwable {
        makePutRequestWithPayload(LINK_RULE_URL, PAYLOAD, LINK_RULE_KEY);
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
        delete(LINK_RULE_URL, LINK_RULE_KEY);
    }

    @Given("^the link rule URL of the service with no query parameters exist$")
    public void theLinkRuleURLOfTheServiceWithNoQueryParametersExist() {
        util.put(ENDPOINT_URI, String.format("%s/linkrule", baseUrl));
    }

    @When("^a request to GET a list of link rules is done$")
    public void aRequestToGETAListOfLinkRulesIsDone() {
        get(LinkRuleList.class);
    }

    @Given("^an inventory file exists$")
    public void anInventoryFileExists() {
        util.put(INVENTORY_FILE, util.getFile("inventory/ansible_hosts"));
    }

    @Given("^the inventory key is known$")
    public void theInventoryKeyIsKnown() {
        util.put(INVENTORY_KEY, "test_inventory");
    }

    @Given("^the inventory upload URL is known$")
    public void theInventoryUploadURLIsKnown() {
        util.put(INVENTORY_URL, String.format("%s/inventory/{key}", baseUrl));
    }

    @When("^an HTTP PUT request with the inventory payload is executed$")
    public void anHTTPPUTRequestWithTheInventoryPayloadIsExecuted() {
        String inventory = util.get(INVENTORY_FILE);
        String url = util.get(INVENTORY_URL);
        String key = util.get(INVENTORY_KEY);
        Map<String, Object> vars = new HashMap<>();
        vars.put("key", key);
        ResponseEntity<Result> response = null;
        try {
            response = client.exchange(url, HttpMethod.PUT, getEntity(inventory), Result.class, vars);
            util.put(RESPONSE, response);
            util.remove(EXCEPTION);
        }
        catch (Exception ex) {
            util.put(EXCEPTION, ex);
        }
    }

    @Then("^the inventory config item is created$")
    public void theInventoryConfigItemIsCreated() {

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
    private void delete(String urlKeyLabel, String resourceKeyLabel) {
        String url = (String) util.get(urlKeyLabel);
        ResponseEntity<Result> response = null;
        try {
            if (resourceKeyLabel != null) {
                response = client.exchange(url, HttpMethod.DELETE, null, Result.class, (String) util.get(resourceKeyLabel));
            } else {
                response = client.exchange(url, HttpMethod.DELETE, null, Result.class);
            }
            util.put(RESPONSE, response);
            util.remove(EXCEPTION);
        }
        catch (Exception ex) {
            util.put(EXCEPTION, ex);
        }
    }

    /*
        a generic get to an endpoint without parameters
     */
    private <T> void get(Class<T> cls){
        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
        ResponseEntity<T> result = client.exchange(
                (String)util.get(ENDPOINT_URI),
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

    @Given("^the inventory exists in the database$")
    public void theInventoryExistsInTheDatabase() {
        theInventoryKeyIsKnown();
        anInventoryFileExists();
        theInventoryUploadURLIsKnown();
        anHTTPPUTRequestWithTheInventoryPayloadIsExecuted();
    }

    @Given("^the URL of the inventory finder endpoint is known$")
    public void theURLOfTheInventoryFinderEndpointIsKnown() {
        util.put(INVENTORY_URL, String.format("%sinventory/{key}", baseUrl));
    }

    @When("^an HTTP GET to the inventory GET endpoint is made using its key$")
    public void anHTTPGETToTheInventoryGETEndpointIsMadeUsingItsKey() {
        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
        ResponseEntity<LinkList> result = client.exchange(
                util.get(INVENTORY_URL),
                HttpMethod.GET,
                new HttpEntity<>(null, headers),
                LinkList.class,
                (String)util.get(INVENTORY_KEY));
        util.put(RESPONSE, result);
    }

    @Given("^there are items linked to the root item in the database$")
    public void thereAreItemsLinkedToTheRootItemInTheDatabase() {
        // imports an ansible inventory
        anInventoryFileExists();
        theInventoryKeyIsKnown();
        theInventoryUploadURLIsKnown();
        anHTTPPUTRequestWithTheInventoryPayloadIsExecuted();
        thereIsNotAnyErrorInTheResponse();
    }

    @Given("^the URL of the snapshot create endpoint is known$")
    public void theURLOfTheSnapshotCreateEndpointIsKnown() {
        util.put(SNAPSHOT_CREATE_URL, String.format("%s/snapshot", baseUrl));
    }

    @When("^a snapshot creation is requested$")
    public void aSnapshotCreationIsRequested() {
        String url = util.get(SNAPSHOT_CREATE_URL);
        String payload = util.get(SNAPSHOT_CREATE_PAYLOAD);
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

    @Given("^a payload exists with the data required to create the snapshot$")
    public void aPayloadExistsWithTheDataRequiredToCreateTheSnapshot() {
        util.put(SNAPSHOT_CREATE_PAYLOAD, util.getFile("payload/create_snapshot.json"));
    }

    @Then("^the result contains no errors$")
    public void theResultContainsNoErrors() {
        ResponseEntity<Result> result = util.get(RESPONSE);
        if(result.getBody().isError()){
            throw new RuntimeException(String.format("Result contains an error as follows: '%s'", result.getBody().getMessage()));
        };
    }

    @Given("^the URL of the snapshot update endpoint is known$")
    public void theURLOfTheSnapshotUpdateEndpointIsKnown() {
        util.put(SNAPSHOT_UPDATE_URL, String.format("%s/snapshot/{root_item_key}/{label}", baseUrl));
    }

    @Given("^the snapshot already exists$")
    public void theSnapshotAlreadyExists() {
        try {
            // replays create_snapshot.feature
            theURLOfTheSnapshotCreateEndpointIsKnown();
            thereAreItemsLinkedToTheRootItemInTheDatabase();
            aPayloadExistsWithTheDataRequiredToCreateTheSnapshot();
            aSnapshotCreationIsRequested();
            theResponseCodeIs(200);
            theResultContainsNoErrors();
        } catch (Exception ex) {
            if (ex.getMessage().contains("duplicate key value violates unique constraint")) {
                // the snapshot is already in the database so do nothing
            } else {
                throw ex;
            }
        }
    }

    @Given("^a payload exists with the data required to update the snapshot$")
    public void aPayloadExistsWithTheDataRequiredToUpdateTheSnapshot() {
        util.put(SNAPSHOT_UPDATE_PAYLOAD, util.getFile("payload/update_snapshot.json"));
    }

    @When("^a snapshot update is requested$")
    public void aSnapshotUpdateIsRequested() {
        String url = util.get(SNAPSHOT_UPDATE_URL);
        String payload = util.get(SNAPSHOT_UPDATE_PAYLOAD);
        String currentLabel = util.get(SNAPSHOT_LABEL);
        String itemRootKey = util.get(SNAPSHOT_ITEM_ROOT_KEY);
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

    @Given("^the item root key of the snapshot is known$")
    public void theItemRootKeyOfTheSnapshotIsKnown() {
        util.put(SNAPSHOT_ITEM_ROOT_KEY, "test_inventory");
    }

    @Given("^the current label of the snapshot is known$")
    public void theCurrentLabelOfTheSnapshotIsKnown() {
        util.put(SNAPSHOT_LABEL, "v1");
    }

    @Given("^the URL of the snapshot delete endpoint is known$")
    public void theURLOfTheSnapshotDeleteEndpointIsKnown() {
        util.put(SNAPSHOT_DELETE_URL, String.format("%s/snapshot/{root_item_key}/{label}", baseUrl));
    }

    @When("^a snapshot delete is requested$")
    public void aSnapshotDeleteIsRequested() {
        String url = util.get(SNAPSHOT_DELETE_URL);
        String currentLabel = util.get(SNAPSHOT_LABEL);
        String itemRootKey = util.get(SNAPSHOT_ITEM_ROOT_KEY);
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

    @Given("^the URL of the snapshot get endpoint is known$")
    public void theURLOfTheSnapshotGetEndpointIsKnown() {
        util.put(SNAPSHOT_LIST_URL, String.format("%s/snapshot/{root_item_key}", baseUrl));
    }

    @Given("^there are snapshots for a given item in the database$")
    public void thereAreSnapshotsForAGivenItemInTheDatabase() {
        // for now only puts one snapshot in the database
        theSnapshotAlreadyExists();
    }

    @Given("^the item root key of the snapshots is known$")
    public void theItemRootKeyOfTheSnapshotsIsKnown() {
        util.put(SNAPSHOT_ITEM_ROOT_KEY, "test_inventory");
    }

    @When("^a snapshot list for an item is requested$")
    public void aSnapshotListForAnItemIsRequested() {
        String url = util.get(SNAPSHOT_LIST_URL);
        String itemRootKey = util.get(SNAPSHOT_ITEM_ROOT_KEY);
        Map<String, Object> vars = new HashMap<>();
        vars.put("root_item_key", itemRootKey);
        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
        ResponseEntity<SnapshotList> result = client.exchange(
                url,
                HttpMethod.GET,
                new HttpEntity<>(null, headers),
                SnapshotList.class,
                vars);
        util.put(RESPONSE, result);
    }

    @Then("^the response contains more than (\\d+) snapshots$")
    public void theResponseContainsMoreThanSnapshots(int count) {
        ResponseEntity<SnapshotList> response = util.get(RESPONSE);

        SnapshotList items = response.getBody();
        if (items != null) {
            if (items.getItems().size() <= count) {
                throw new RuntimeException(
                    String.format(
                        "Response does not contain more than '%s' but '%s' snapshots.",
                        count,
                        response.getBody().getItems().size()
                    )
                );
            }
        }
        else {
            throw new RuntimeException(
                String.format(
                    "Response contains no snapshots where more than '%s' were expected.",
                    count
                )
            );
        }
    }

    @Given("^the snapshot does not already exist$")
    public void theSnapshotDoesNotAlreadyExist() {
        // remove the snapshot if exists
        theURLOfTheSnapshotDeleteEndpointIsKnown();
        theItemRootKeyOfTheSnapshotIsKnown();
        theCurrentLabelOfTheSnapshotIsKnown();
        aSnapshotDeleteIsRequested();
        theResponseCodeIs(200);
    }

    @Given("^the URL of the item tree get endpoint is known$")
    public void theURLOfTheItemTreeGetEndpointIsKnown() {
        util.put(SNAPSHOT_TREE_URL, String.format("%s/tree/{root_item_key}/{label}", baseUrl));
    }

    @When("^a snapshot tree retrieval for the snapshot is requested$")
    public void aSnapshotTreeRetrievalForTheSnapshotIsRequested() {
        String url = util.get(SNAPSHOT_TREE_URL);
        String itemRootKey = util.get(SNAPSHOT_ITEM_ROOT_KEY);
        String label = util.get(SNAPSHOT_LABEL);
        Map<String, Object> vars = new HashMap<>();
        vars.put("root_item_key", itemRootKey);
        vars.put("label", label);
        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
        ResponseEntity<ItemTreeData> result = client.exchange(
                url,
                HttpMethod.GET,
                new HttpEntity<>(null, headers),
                ItemTreeData.class,
                vars);
        util.put(RESPONSE, result);
    }

    @Then("^the result contains the tree items and links$")
    public void theResultContainsTheTreeItemsAndLinks() {
        ResponseEntity<ItemTreeData> response = util.get(RESPONSE);
        ItemTreeData tree = response.getBody();
        if (tree != null) {
            if (tree.getItems().size() == 0 || tree.getLinks().size() == 0){
                throw new RuntimeException("Tree does not contain items or links.");
            }
        }
    }
}
