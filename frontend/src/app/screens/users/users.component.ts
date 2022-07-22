import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { faUsers } from '@fortawesome/free-solid-svg-icons';

@Component({
  selector: 'app-users',
  templateUrl: './users.component.html',
  styleUrls: ['./users.component.scss'],
})
export class UsersComponent {
  public faUsers = faUsers;
  public newUserForm: FormGroup;

  constructor(private fb: FormBuilder) {
    this.newUserForm = this.buildForm();
  }

  private buildForm(): FormGroup {
    return this.fb.group({
      name: ['', Validators.required],
      password: ['', Validators.required],
      isAdmin: [false],
    });
  }

  public onSubmit() {
    console.log('form', this.newUserForm.value);
  }
}
