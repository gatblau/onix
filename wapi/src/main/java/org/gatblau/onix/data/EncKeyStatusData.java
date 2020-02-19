package org.gatblau.onix.data;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

import java.io.Serializable;

@ApiModel
public class EncKeyStatusData implements Serializable {
    private static final long serialVersionUID = 1L;
    private long key1;
    private long key2;

    @ApiModelProperty(
            position = 1,
            required = true,
            value = "The number of items using Key 1.")
    public long getKey1() {
        return key1;
    }

    public void setKey1(String key1) {
        this.key1 = Long.parseLong(key1);
    }

    @ApiModelProperty(
            position = 1,
            required = true,
            value = "The number of items using Key 2.")
    public long getKey2() {
        return key2;
    }

    public void setKey2(String key2) {
        this.key2 = Long.parseLong(key2);
    }

    @ApiModelProperty(
            position = 1,
            required = true,
            value = "The percentage of items using key 1.")
    public double getKey1PercentageUse() {
        try {
            double v = (double)key1 / (double)(key1 + key2);
            return v * 100;
        } catch (Exception e) {
        }
        return 0;
    }

    @ApiModelProperty(
            position = 1,
            required = true,
            value = "The percentage of items using key 2.")
    public double getKey2PercentageUse() {
        try {
            double v = (double)key2 / (double)(key1 + key2);
            return v * 100;
        } catch (Exception e) {
        }
        return 0;
    }
}
