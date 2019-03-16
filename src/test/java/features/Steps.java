package features;

import cucumber.api.java.en.And;
import cucumber.api.java.en.Given;
import cucumber.api.java.en.Then;
import cucumber.api.java.en.When;
import org.gatblau.onix.Info;
import org.gatblau.onix.data.*;
import org.json.simple.JSONObject;
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
        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "text/html");
        ResponseEntity<String> result = client.exchange(
                (String)util.get(BASE_URL),
                HttpMethod.GET,
                new HttpEntity<>(null, headers),
                String.class);
        util.put(RESPONSE, result);
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
        util.put(CONGIG_ITEM_TYPE_KEY, "ANSIBLE_HOST");
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

    @Then("^the response contains more than (\\d+) items$")
    public void theResponseContainsMoreThanNumberItems(int count) {
        ResponseEntity<ItemList> response = util.get(RESPONSE);

        ItemList items = response.getBody();
        if (items != null) {
            if (items.getItems().size() <= count) {
                throw new RuntimeException(
                    String.format(
                        "Response does not contain more than '%s' items but '%s' items.",
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

    @Then("^the response contains more than (\\d+) item types$")
    public void theResponseContainsMoreThanItemTypes(int items) {
        ResponseEntity<ItemTypeList> response = util.get(RESPONSE);
        int actual = response.getBody().getItems().size();
        if(response.getBody().getItems().size() <= items){
            throw new RuntimeException(String.format("Response contains %s items instead of %s items.", actual, items));
        }
    }

    @Then("^the response contains more than (\\d+) link rules$")
    public void theResponseContainsMoreThanLinkRules(int rules) {
        ResponseEntity<LinkRuleList> response = util.get(RESPONSE);
        int actual = response.getBody().getItems().size();
        if(response.getBody().getItems().size() <= rules){
            throw new RuntimeException(String.format("Response contains %s items which is less than %s items.", actual, rules));
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
        theLinkTypeNaturalKeyIsKnown();
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
        try {
            // replays create_tag.feature
            theURLOfTheTagCreateEndpointIsKnown();
            putData("payload/create_data_payload.json");
            aPayloadExistsWithTheDataRequiredToCreateTheTag();
            aTagCreationIsRequested();
            theResponseCodeIs(200);
            theResultContainsNoErrors();
        } catch (Exception ex) {
            if (ex.getMessage().contains("duplicate key value violates unique constraint")) {
                // the tag is already in the database so do nothing
            } else {
                throw ex;
            }
        }
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
            if (items.getItems().size() <= count) {
                throw new RuntimeException(
                    String.format(
                        "Response does not contain more than '%s' but '%s' tags.",
                        count,
                        response.getBody().getItems().size()
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
        util.put(TAG_TREE_URL, String.format("%s/tree/{root_item_key}/{label}", baseUrl));
    }

    @When("^a tag tree retrieval for the tag is requested$")
    public void aTagTreeRetrievalForTheTagIsRequested() {
        String url = util.get(TAG_TREE_URL);
        String itemRootKey = util.get(ROOT_ITEM_KEY);
        String label = util.get(TAG_LABEL);
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
        for (Result result : results.getBody().getItems()) {
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
        for (Result result : results.getItems()) {
            updateFound = result.getOperation().equals("U");
            if (updateFound) break;
        }
        if (!updateFound){
            throw new RuntimeException("Results have not been updated");
        }
    }

    @Given("^the URL of the item tree DELETE endpoint is known$")
    public void theURLOfTheItemTreeDELETEEndpointIsKnown() {
        util.put(DELETE_TREE_URL, String.format("%s/tree/{root_item_key}", baseUrl));
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

    @Then("^the reponse contains the requested metadata$")
    public void theReponseContainsTheRequestedMetadata() {
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
        theItemTypeURLOfTheServiceWithKeyIsKnown();
        theItemTypeNaturalKeyIsKnown();
        aJsonPayloadWithNewItemTypeInformationExists();
        anItemTypePUTHTTPRequestWithAJSONPayloadIsDone();
    }

    @Given("^the item with metadata exists in the database$")
    public void theItemWithMetadataExistsInTheDatabase() {
        util.put(ITEM_KEY, ITEM_ONE_KEY);
        putItem(util.get(ITEM_KEY), "payload/create_meta_test_item_payload.json");
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

    private void putData(String payloadFilePath){
        String payload = util.getFile(payloadFilePath);
        String url = String.format("%s/data", baseUrl);
        Map<String, Object> vars = new HashMap<>();
        vars.put("payload", util.getFile(payloadFilePath));
        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
        ResponseEntity<ResultList> response = client.exchange(url, HttpMethod.PUT, getEntity(payload), ResultList.class, vars);
    }

    @Given("^link rules exist in the database$")
    public void linkRulesExistInTheDatabase() {
        putData("payload/import_link_rules_payload.json");
    }
}
