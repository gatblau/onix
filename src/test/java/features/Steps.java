package features;

import cucumber.api.PendingException;
import cucumber.api.java.en.And;
import cucumber.api.java.en.Given;
import cucumber.api.java.en.When;
import org.gatblau.onix.Info;
import org.gatblau.onix.Result;
import org.springframework.http.*;

import static features.Key.*;

import javax.annotation.PostConstruct;
import java.util.HashMap;
import java.util.Map;

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

    @And("^the item URL of the service is known$")
    public void theItemURLOfTheServiceIsKnown() throws Throwable {
        util.put(Key.ITEM_URL, String.format("%s/item/{key}/", baseUrl));
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

    @And("^a yaml payload with node information exists$")
    public void aYamlPayloadWithNodeInformationExists() throws Throwable {
        String payload = util.getFile("payload/create_item_payload.yml");
        util.put(Key.PAYLOAD, payload);
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
        String payload = util.get(PAYLOAD);
        String url = util.get(ITEM_URL);
        Map<String, Object> vars = new HashMap<>();
        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
        HttpEntity<?> entity = new HttpEntity<>(payload, headers);
        vars.put("key", util.get(ITEM_ONE_KEY));
        ResponseEntity<Result> response = null;
        try {
            response = client.exchange(url, HttpMethod.PUT, entity, Result.class, vars);
            util.put(RESPONSE, response);
            util.remove(EXCEPTION);
        }
        catch (Exception ex) {
            util.put(EXCEPTION, ex);
        }
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
        theItemURLOfTheServiceIsKnown();
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

    @Given("^the natural keys for two configuration items are known$")
    public void theNaturalKeysForTwoConfigurationItemsAreKnown() throws Throwable {
        util.put(ITEM_ONE_KEY, "ITEM_ONE_KEY");
        util.put(ITEM_TWO_KEY, "ITEM_TWO_KEY");
    }

    @Given("^the link URL of the service is known$")
    public void theLinkURLOfTheServiceIsKnown() throws Throwable {
        util.put(LINK_URL, String.format("%s/link/{fromItemKey}/{toItemKey}/", baseUrl));
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
            client.delete((String) util.get(LINK_URL), (String)util.get(ITEM_ONE_KEY), (String)util.get(ITEM_TWO_KEY));
            util.remove(EXCEPTION);
        }
        catch (Exception ex) {
            util.put(EXCEPTION, ex);
        }
    }
}
