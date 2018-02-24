package org.gatblau.onix.model;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;

import javax.persistence.AttributeConverter;
import javax.persistence.Converter;
import java.io.IOException;

@Converter(autoApply = true)
public class JSONBConverter implements AttributeConverter<JsonNode, String> {
    private static ObjectMapper mapper = new ObjectMapper();

    static {
        mapper.setSerializationInclusion(JsonInclude.Include.NON_NULL);
    }

    @Override
    public String convertToDatabaseColumn(final JsonNode objectValue) {
        return (objectValue == null) ? null : objectValue.toString();
    }

    @Override
    public JsonNode convertToEntityAttribute(final String dataValue) {
        try {
            return (dataValue == null) ? null : mapper.readTree(dataValue);
        } catch (IOException e) {
            throw new RuntimeException("Unable to deserialize to json field", e);
        }
    }
}
