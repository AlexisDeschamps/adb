<template>
  <div class="address-wrapper">
    <!-- address 1 input needs to have type "search" to prevent
         chrome from autofilling addresses -->
    <input
      id="canvassingAddress1"
      type="search"
      class="form-control address-input-margin-bottom"
      placeholder="Address 1"
      v-model="address1"
      @change="emitChangeEvent"
      @keyup="emitChangeEvent"
      />
    <input
      type="text"
      class="form-control address-input-margin-bottom"
      placeholder="Address 2"
      v-model="address2"
      @change="emitChangeEvent"
      @keyup="emitChangeEvent"
      />
    <input
      type="text"
      class="form-control address-input-margin-bottom"
      placeholder="City"
      v-model="city"
      @change="emitChangeEvent"
      @keyup="emitChangeEvent"
      />
    <input
      type="text"
      class="form-control address-input-margin-bottom"
      placeholder="State"
      v-model="state"
      @change="emitChangeEvent"
      @keyup="emitChangeEvent"
      />
    <input
      type="text"
      class="form-control"
      placeholder="ZIP code"
      v-model="zip"
      @change="emitChangeEvent"
      @keyup="emitChangeEvent"
      />
  </div>
</template>

<script>
import Vue from 'vue';

declare global {
  const google: any;
}

export default Vue.extend({
  props: {
    initialAddress1: String,
    initialAddress2: String,
    initialCity: String,
    initialState: String,
    initialZip: String,
  },
  data() {
    return {
      /**
       * The Autocomplete object.
       *
       * @type {Autocomplete}
       * @link https://developers.google.com/maps/documentation/javascript/reference#Autocomplete
       */
      autocomplete: null as any,

      address1: '',
      address2: '',
      city: '',
      state: '',
      zip: '',
    };
  },

  mounted: function() {
    const options: any = {};
    options.componentRestrictions = {
      country: "us",
    };

    this.autocomplete = new (google as any).maps.places.Autocomplete(
      document.getElementById("canvassingAddress1"),
      options
    );
    this.autocomplete.addListener('place_changed', this.onPlaceChanged);
  },

  methods: {

    emitChangeEvent() {
      this.$emit('change', {
        address1: this.address1,
        address2: this.address2,
        city: this.city,
        state: this.state,
        zip: this.zip,
      });
    },

    onPlaceChanged() {
      let place = this.autocomplete.getPlace();

      if (!place.geometry) {
        console.log("no place found");
        return;
      }

      if (place.address_components !== undefined) {
        // reset address
        this.zip = this.state = this.city = this.address2 = this.address1 = '';

        let streetNum = '';
        let route = '';
        for (let addressComponent of place.address_components) {
          for (let compType of addressComponent.types) {
            if (compType === "street_number") {
              streetNum = addressComponent.short_name;
            } else if (compType === "route") {
              route = addressComponent.short_name;
            } else if (compType === "locality") {
              this.city = addressComponent.short_name;
            } else if (compType === "administrative_area_level_1") {
              this.state = addressComponent.short_name;
            } else if (compType === "postal_code") {
              this.zip = addressComponent.short_name;
            }
          }
        }
        this.address1 = streetNum + ' ' + route;
        this.emitChangeEvent();
      }
    },
  },
  watch: {
    initialAddress1: function() {
      this.address1 = this.initialAddress1;
    },
    initialAddress2: function() {
      this.address2 = this.initialAddress2;
    },
    initialCity: function() {
      this.city = this.initialCity;
    },
    initialState: function() {
      this.state = this.initialState;
    },
    initialZip: function() {
      this.zip = this.initialZip;
    },
  }
})
</script>

<style>
  .address-input-margin-bottom {
    margin-bottom: 4px;
  }
</style>
