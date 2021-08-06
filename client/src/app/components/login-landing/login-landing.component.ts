import { Component, OnInit } from '@angular/core';
import { FormBuilder, Validators } from '@angular/forms';
import { LoginService } from 'src/app/services/login-service/login.service'
import { LoginResponse } from 'src/app/services/loginbackend-service/loginbackend.service'
import { Router } from '@angular/router';

class LoginInfo {
  name:string;
  id:string;
  constructor(name:string, id:string) {
    this.name = name;
    this.id = id;
  };
};

@Component({
  selector: 'app-login-landing',
  templateUrl: './login-landing.component.html',
  styleUrls: ['./login-landing.component.less'],
  providers: [ ]
})

export class LoginLandingComponent implements OnInit {

  loginForm = this.formBuilder.group({
    name: ['', Validators.required],
  });

  constructor(private formBuilder: FormBuilder, private loginService: LoginService, private router: Router) { 

  }
  
  onSubmit() {
    const name = this.loginForm.value.name;
    this.loginService.login(name).subscribe({
      next: (resp: LoginResponse)=> {
        if (resp.success) {
          this.router.navigate(["/main"]);
        }
      }
    });
  }

  ngOnInit(): void {
  
  }

}
