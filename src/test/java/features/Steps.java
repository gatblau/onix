package features;

import cucumber.api.PendingException;
import cucumber.api.java.en.And;
import cucumber.api.java.en.Then;
import cucumber.api.java.en.When;
import org.gatblau.onix.Info;
import org.gatblau.onix.Result;
import org.springframework.http.*;

import static features.Key.*;

import javax.annotation.PostConstruct;
import java.io.File;
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

    @And("^the create URL of the service is known$")
    public void theCreateURLOfTheServiceIsKnown() throws Throwable {
        util.put(Key.CREATE_URL, String.format("%s/item/{key}/", baseUrl));
    }

    @And("^the service responds with the create node Id$")
    public void theServiceRespondsWithTheCreateNodeId() throws Throwable {
        ResponseEntity<Result> response = util.get(RESPONSE);
        assert (Integer.parseInt(response.getBody().getValue()) > 0);
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

    @And("^a json payload with node information exists$")
    public void aJsonPayloadWithNodeInformationExists() throws Throwable {
        String payload = util.getFile("payload/create_item_payload.json");
        util.put(Key.PAYLOAD, payload);
    }

    @And("^the node does not exist in the database$")
    public void theNodeDoesNotExistInTheDatabase() throws Throwable {
        theClearURLOfTheServiceIsKnown();
        aDeleteRequestToTheServiceIsDone();
    }

    @And("^the database is cleared$")
    public void theDatabaseIsCleared() throws Throwable {
        assert(!util.containsKey(EXCEPTION));
    }

    @And("^the clear URL of the service is known$")
    public void theClearURLOfTheServiceIsKnown() throws Throwable {
        util.put(Key.CLEAR_URL, String.format("%s/clear/", baseUrl));
    }

    @And("^a delete request to the service is done$")
    public void aDeleteRequestToTheServiceIsDone() throws Throwable {
        try {
            client.delete((String) util.get(CLEAR_URL));
            util.remove(EXCEPTION);
        }
        catch (Exception ex) {
            util.put(EXCEPTION, ex);
        }
    }

    @Then("^there is not any error in the response$")
    public void thereIsNotAnyErrorInTheResponse() throws Throwable {
        assert(!util.containsKey(EXCEPTION));
    }

    @And("^a node key is known$")
    public void aNodeKeyIsKnown() throws Throwable {
        util.put(ITEM_KEY, "Test_Item_1");
    }

    @When("^a create JSON request to the service is done$")
    public void aCreateJSONRequestToTheServiceIsDone() throws Throwable {
        String payload = util.get(PAYLOAD);
        String url = util.get(CREATE_URL);
        Map<String, Object> vars = new HashMap<>();
        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
        HttpEntity<?> entity = new HttpEntity<>(payload, headers);
        vars.put("key", util.get(ITEM_KEY));
        client.exchange(url, HttpMethod.PUT, entity, Void.class, vars);
    }
}
