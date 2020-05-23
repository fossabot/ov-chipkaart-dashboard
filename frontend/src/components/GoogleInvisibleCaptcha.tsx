import Anchor from '@material-ui/core/Link';
import Typography from '@material-ui/core/Typography';
import React from 'react';
import { useTranslation } from 'react-i18next';
import { createStyles, makeStyles, Theme } from '@material-ui/core/styles';
import ReCAPTCHA from 'react-google-recaptcha';

const useStyles = makeStyles((theme: Theme) =>
    createStyles({
        privacyPolicy: {
            color: theme.palette.grey[700],
        },
    }),
);

export default function GoogleInvisibleCaptcha() {
    const classes = useStyles();
    const { t } = useTranslation();
    return (
        <div>
            <Typography variant="body2" className={classes.privacyPolicy}>
                {t(
                    'This page is protected by reCAPTCHA and is subject to the Google',
                )}{' '}
                <Anchor
                    target="_blank"
                    href="https://policies.google.com/privacy"
                >
                    {t('Privacy Policy')}
                </Anchor>{' '}
                and{' '}
                <Anchor
                    target="_blank"
                    href="https://policies.google.com/terms"
                >
                    {t('Terms of Service')}
                </Anchor>
                .
            </Typography>
            <ReCAPTCHA
                sitekey="6LdGNPsUAAAAAH-De07bVDCXyqX70pPMOyEE7rND"
                size="invisible"
            />
        </div>
    );
}
