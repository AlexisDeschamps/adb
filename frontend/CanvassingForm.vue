<template>
  <adb-page
    :title="'Phone Bank Form'">
    <form action id="canvassingForm" autocomplete="off">
      <label for="firstName">
        First Name
      </label>
      <input id="firstName" class="form-control" v-model="firstName" /> <br />
      <label for="lastName">
        Last Name
      </label>
      <input id="lastName" class="form-control" v-model="lastName" /> <br />
      <label for="email">
        Email
      </label>
      <input id="email" class="form-control" v-model="email" /> <br />
      <label for="phone">
        Phone
      </label>
      <input id="phone" class="form-control" v-model="phone" /> <br />

      <label>Location</label>
      <canvassing-address
        v-on:change="onAddressChange"
        /> <br />

      <label for="source">Source</label>
      <select id="source" class="form-control" v-model="source">
        <option value="phone">Phone Banking</option>
        <option value="canvass">In Person Canvassing</option>
        <option value="canvass">House Party</option>
      </select> <br />
      <label for="dateSourced">Date Sourced</label>
      <input id="dateSourced" class="form-control" type="date" v-model="dateSourced" /> <br />

      <input type="checkbox" id="voter" v-model="voter" />
      <label for="voter"> Eligible Berkeley Voter </label> <br />

      <div>
        <h3> Requests </h3>
        <input type="checkbox" id="requestedLawnSign" v-model="requestedLawnSign" />
        <label for="requestedLawnSign"> Requested Lawn Sign </label> <br />
        <input type="checkbox" id="requestedPoster" v-model="requestedPoster" />
        <label for="requestedPoster"> Requested Poster </label> <br />
        <br />
      </div>

      <div>
        <h3>Issues</h3>
        <input type="checkbox" id="issueHousing" v-model="issueHousing" />
        <label for="issueHousing"> Housing </label> <br />
        <input type="checkbox" id="issueHomelessness" v-model="issueHomelessness" />
        <label for="issueHomelessness"> Homelessness </label> <br />
        <input type="checkbox" id="issueClimate" v-model="issueClimate" />
        <label for="issueClimate"> Climate </label> <br />
        <input type="checkbox" id="issuePublicSafety" v-model="issuePublicSafety" />
        <label for="issuePublicSafety"> Public Safety </label> <br />
        <input type="checkbox" id="issuePoliceAccountability" v-model="issuePoliceAccountability" />
        <label for="issuePoliceAccountability"> Police Accountability </label> <br />
        <input type="checkbox" id="issueTransit" v-model="issueTransit" />
        <label for="issueTransit"> Transit </label> <br />
        <input type="checkbox" id="issueEconomicEquality" v-model="issueEconomicEquality" />
        <label for="issueEconomicEquality"> Economic Equality </label> <br />
        <input type="checkbox" id="issuePublicHealth" v-model="issuePublicHealth" />
        <label for="issuePublicHealth"> Public Health </label> <br />
        <input type="checkbox" id="issueAnimalRights" v-model="issueAnimalRights" />
        <label for="issueAnimalRights"> Animal Rights </label> <br />

        <br />
      </div>

      <div>
        <h3>Support</h3>

        <input type="checkbox" id="interestDonate" v-model="interestDonate" />
        <label for="interestDonate"> Donating </label> <br />

        <input type="checkbox" id="interestAttendEvent" v-model="interestAttendEvent" />
        <label for="interestAttendEvent"> Attend Event </label> <br />

        <input type="checkbox" id="interestVolunteer" v-model="interestVolunteer" />
        <label for="interestVolunteer"> Volunteer </label> <br />

        <input type="checkbox" id="interestHostEvent" v-model="interestHostEvent" />
        <label for="interestHostEvent"> Host Event (e.g. house party) </label> <br />

        <br />
      </div>

      <input type="checkbox" id="requiresFollowup" v-model="requiresFollowup" />
      <label for="requiresFollowup"> Requires Followup (for something not covered by another option) </label> <br />

      <br />
      <label for="notes"> Notes </label>
      <textarea id="notes" type="text" class="form-control" v-model="notes" /> <br />

      <!--
      <label for="canvasser"> Canvasser </label>
      <input id="canvasser" type="text" class="form-control" v-model="canvasser" /> <br />

      <label for="canvassLeader"> Canvass Leader </label>
      <input id="canvassLeader" type="text" class="form-control" v-model="canvassLeader" /> <br />
      -->

    </form>
    <center>
      <button
        class="btn btn-success btn-lg"
        id="submit-button"
        v-on:click="save"
        :disabled="saving"
        >
        <span>Save</span>
      </button>
    </center>
  </adb-page>


</template>

<script lang="ts">
import Vue from 'vue';
import AdbPage from './AdbPage.vue';
import CanvassingAddress from './components/CanvassingAddress.vue';
import { flashMessage } from './flash_message';

function getDateTodayStr() {
  var d = new Date();
  var year = '' + d.getFullYear();
  var rawMonth = d.getMonth() + 1;
  var month = rawMonth > 9 ? '' + rawMonth : '0' + rawMonth;
  var rawDate = d.getDate();
  var date = rawDate > 9 ? '' + rawDate : '0' + rawDate;
  var validDateString = year + '-' + month + '-' + date;
  return validDateString;
}

export default Vue.extend({
  components: {
    AdbPage,
    CanvassingAddress,
  },
  data() {
    return {
      saving: false,

      // Supporter fields:
      firstName: '',
      lastName: '',
      email: '',
      phone: '',
      address1: '',
      address2: '',
      city: '',
      state: '',
      zip: '',
      source: 'phone',
      dateSourced: getDateTodayStr(),
      requestedLawnSign: false,
      requestedPoster: false,
      voter: true,

      issueHousing: false,
      issueHomelessness: false,
      issueClimate: false,
      issuePublicSafety: false,
      issuePoliceAccountability: false,
      issueTransit: false,
      issueEconomicEquality: false,
      issuePublicHealth: false,
      issueAnimalRights: false,

      interestDonate: false,
      interestAttendEvent: false,
      interestVolunteer: false,
      interestHostEvent: false,

      requiresFollowup: false,

      notes: '',

      // canvasser: 'Samer Masterson',
      // canvassLeader: 'Samer Masterson',
    }
  },
  methods: {
    // TODO: use correct type for 'e'
    onAddressChange(e: any) {
      console.log(e);
      this.address1 = e.address1;
      this.address2 = e.address2;
      this.city = e.city;
      this.state = e.state;
      this.zip = e.zip;
    },

    save() {
      if (this.phone.trim() === '' && this.email.trim() === '') {
        flashMessage("Error: At least one of phone or email must be set.", true);
        return;
      }
      this.saving = true;
      $.ajax({
        url: '/canvass/supporter/save',
        method: 'POST',
        contentType: 'application/json',
        data: JSON.stringify({
          first_name: this.firstName.trim(),
          last_name: this.lastName.trim(),
          email: this.email.trim(),
          phone: this.phone.trim(),
          location_address1: this.address1.trim(),
          location_address2: this.address2.trim(),
          location_city: this.city.trim(),
          location_state: this.state.trim(),
          location_zip: this.zip.trim(),
          source: this.source,
          date_sourced: this.dateSourced.trim(),
          requested_lawn_sign: this.requestedLawnSign,
          requested_poster: this.requestedPoster,
          voter: this.voter,
          issue_housing: this.issueHousing,
          issue_homelessness: this.issueHomelessness,
          issue_climate: this.issueClimate,
          issue_public_safety: this.issuePublicSafety,
          issue_police_accountability: this.issuePoliceAccountability,
          issue_transit: this.issueTransit,
          issue_economic_equality: this.issueEconomicEquality,
          issue_public_health: this.issuePublicHealth,
          issue_animal_rights: this.issueAnimalRights,
          interest_donate: this.interestDonate,
          interest_attend_event: this.interestAttendEvent,
          interest_volunteer: this.interestVolunteer,
          interest_host_event: this.interestHostEvent,
          requires_followup: this.requiresFollowup,
          notes: this.notes.trim(),
          // canvasser: this.canvasser,
          // canvass_leader: this.canvassLeader,
        }),
        success: (data) => {
          this.saving = false;
          let parsed = JSON.parse(data);
          if (parsed.status === 'error') {
            flashMessage('Error: ' + parsed.message, true);
            return;
          }

          flashMessage('Saved!', false);
        },
        error: () => {
          this.saving = false;
          flashMessage("Server error, did not save data.", true);
        },
      })
    },
  }
});

</script>
