package org.gatblau.onix.data;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

import java.io.Serializable;

@ApiModel
public class EncKeyStatusData implements Serializable {
    private static final long serialVersionUID = 1L;
    private long noKeyCount;
    private long key1Count;
    private long key2Count;
    private short activeKey;
    private short defaultKey;
    private String defaultKeyExpiry;

    @ApiModelProperty(
            position = 1,
            required = true,
            value = "The percentage of items using key 2.")
    public long getNoKeyCount() {
        return noKeyCount;
    }

    public void setNoKeyCount(long noKeyCount) {
        this.noKeyCount = noKeyCount;
    }

    @ApiModelProperty(
            position = 2,
            required = true,
            value = "The number of items using Key 1.")
    public long getKey1Count() {
        return key1Count;
    }

    public void setKey1Count(long key1Count) {
        this.key1Count = key1Count;
    }

    @ApiModelProperty(
            position = 3,
            required = true,
            value = "The number of items using Key 2.")
    public long getKey2Count() {
        return key2Count;
    }

    public void setKey2Count(long key2Count) {
        this.key2Count = key2Count;
    }

    @ApiModelProperty(
            position = 4,
            required = true,
            value = "The percentage of items using key 1.")
    public double getKey1v2Ratio() {
        try {
            double v = (double) key1Count / (double)(key1Count + key2Count);
            return v * 100;
        } catch (Exception e) {
        }
        return 0;
    }

    @ApiModelProperty(
            position = 5,
            required = true,
            value = "The percentage of items using key 2.")
    public double getKey2v1Ratio() {
        try {
            double v = (double) key2Count / (double)(key1Count + key2Count);
            return v * 100;
        } catch (Exception e) {
        }
        return 0;
    }

    public short getActiveKey() {
        return activeKey;
    }

    @ApiModelProperty(
            position = 6,
            required = true,
            value = "The key currently in use.")
    public void setActiveKey(short activeKey) {
        this.activeKey = activeKey;
    }

    @ApiModelProperty(
            position = 7,
            required = true,
            value = "The key in use if the expiry date is not passed.")
    public short getDefaultKey() {
        return defaultKey;
    }

    public void setDefaultKey(short defaultKey) {
        this.defaultKey = defaultKey;
    }

    @ApiModelProperty(
            position = 8,
            required = true,
            value = "The expiration date for the default key after which the secondary key is used.")
    public String getDefaultKeyExpiry() {
        return defaultKeyExpiry;
    }

    public void setDefaultKeyExpiry(String defaultKeyExpiry) {
        this.defaultKeyExpiry = defaultKeyExpiry;
    }
}
